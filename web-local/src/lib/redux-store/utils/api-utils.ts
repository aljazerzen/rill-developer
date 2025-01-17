import { createAsyncThunk } from "../redux-toolkit-wrapper";
import type {
  EntityRecord,
  EntityType,
  StateType,
} from "@rilldata/web-local/common/data-modeler-state-service/entity-state-service/EntityStateService";
import type { ActionCreatorWithPreparedPayload } from "@reduxjs/toolkit";
import type { EntityRecordMapType } from "@rilldata/web-local/common/data-modeler-state-service/entity-state-service/EntityStateServicesMap";
import { fetchWrapper } from "../../util/fetchWrapper";
import type { ValidationConfig } from "./validation-utils";
import { validateEntity } from "./validation-utils";
import type { RillReduxState } from "../store-root";

function getQueryArgs(args: Record<string, any>) {
  if (!args) return "";
  return "/?" + Object.keys(args).map((argKey) => `${argKey}=${args[argKey]}`);
}

export function generateApis<
  Type extends EntityType,
  FetchManyParams extends Record<string, any> = Record<string, unknown>,
  CreateParams extends Record<string, any> = Record<string, unknown>,
  Entity extends EntityRecord = EntityRecordMapType[Type][StateType.Persistent]
>(
  [entityType, sliceName, endpoint]: [EntityType, keyof RillReduxState, string],
  [addManyAction, addOneAction, updateAction, removeAction]: [
    ActionCreatorWithPreparedPayload<[entities: Array<Entity>], Array<Entity>>,
    ActionCreatorWithPreparedPayload<[entity: Entity], Entity>,
    ActionCreatorWithPreparedPayload<
      [id: string, changes: Partial<Entity>],
      { id: string; changes: Partial<Entity> }
    >,
    ActionCreatorWithPreparedPayload<[id: string], string>
  ],
  validations: Array<ValidationConfig<Entity>>,
  [createHook, deleteHook]: [
    createHook?: (createdEntity: Entity) => void | Promise<void>,
    deleteHook?: (id: string) => void | Promise<void>
  ] = []
) {
  return {
    // TODO: add caching here to prevent too many fetchManyApi calls
    fetchManyApi: createAsyncThunk(
      `${entityType}/fetchManyApi`,
      async (args: FetchManyParams, thunkAPI) => {
        let entities = await fetchWrapper(
          `${endpoint}${getQueryArgs(args)}`,
          "GET"
        );
        entities = await Promise.all(
          entities.map(async (entity) => {
            const changes = await validateEntity(entity, entity, validations);
            return {
              ...entity,
              ...changes,
            };
          })
        );
        thunkAPI.dispatch(addManyAction(entities));
        return entities;
      }
    ),
    createApi: createAsyncThunk(
      `${entityType}/createApi`,
      async (args: CreateParams, thunkAPI) => {
        const createdEntity = await fetchWrapper(endpoint, "PUT", args);
        thunkAPI.dispatch(addOneAction(createdEntity));
        if (createHook) createHook(createdEntity);
        return createdEntity;
      }
    ),
    updateApi: createAsyncThunk(
      `${entityType}/updateApi`,
      async (
        { id, changes }: { id: string; changes: Partial<Entity> },
        thunkAPI
      ) => {
        const validationChanges = await validateEntity(
          thunkAPI.getState()[sliceName].entities[id] as Entity,
          changes,
          validations
        );
        thunkAPI.dispatch(
          updateAction(
            id,
            await fetchWrapper(`${endpoint}/${id}`, "POST", changes)
          )
        );
        thunkAPI.dispatch(updateAction(id, validationChanges));
      }
    ),
    deleteApi: createAsyncThunk(
      `${entityType}/deleteApi`,
      async (id: string, thunkAPI) => {
        await fetchWrapper(`${endpoint}/${id}`, "DELETE");
        thunkAPI.dispatch(removeAction(id));
        if (deleteHook) deleteHook(id);
      }
    ),
  };
}
