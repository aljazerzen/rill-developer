import type { MeasureDefinitionEntity } from "@rilldata/web-local/common/data-modeler-state-service/entity-state-service/MeasureDefinitionStateService";
import type { ValidationState } from "@rilldata/web-local/common/data-modeler-state-service/entity-state-service/MetricsDefinitionEntityService";
import { metricsExplorerStore } from "../../application-state-stores/explorer-stores";
import { createEntityAdapter, createSlice } from "../redux-toolkit-wrapper";
import type { PayloadAction } from "@reduxjs/toolkit";

const measureDefAdapter = createEntityAdapter<MeasureDefinitionEntity>();

export const measureDefSlice = createSlice({
  name: "measureDefinition",
  initialState: measureDefAdapter.getInitialState(),
  reducers: {
    addManyMeasures: {
      reducer: measureDefAdapter.addMany,
      prepare: (measures: Array<MeasureDefinitionEntity>) => ({
        payload: measures,
      }),
    },

    addOneMeasure: {
      reducer: measureDefAdapter.addOne,
      prepare: (measure: MeasureDefinitionEntity) => ({
        payload: measure,
      }),
    },

    updateMeasure: {
      reducer: measureDefAdapter.updateOne,
      prepare: (id: string, measure: Partial<MeasureDefinitionEntity>) => ({
        payload: { id, changes: measure },
      }),
    },

    removeMeasure: {
      reducer: measureDefAdapter.removeOne,
      prepare: (id: string) => ({ payload: id }),
    },

    clearMeasuresForMetricsDefId: {
      reducer: (state, action: PayloadAction<string>) => {
        measureDefAdapter.removeMany(
          state,
          state.ids.filter(
            (id) => state.entities[id].metricsDefId === action.payload
          )
        );
        metricsExplorerStore.clearLeaderboardMeasureId(action.payload);
      },
      prepare: (id: string) => ({ payload: id }),
    },

    setMeasureExpressionValidation: {
      reducer: (
        state,
        {
          payload: { id, validation, message },
        }: PayloadAction<{
          id: string;
          validation: ValidationState;
          message: string;
        }>
      ) => {
        if (!state.entities[id]) return;
        state.entities[id].expressionIsValid = validation;
        state.entities[id].expressionValidationError = message;
      },
      prepare: (id: string, validation: ValidationState, message: string) => ({
        payload: { id, validation, message },
      }),
    },
  },
});

export const {
  addManyMeasures,
  addOneMeasure,
  updateMeasure,
  removeMeasure,
  clearMeasuresForMetricsDefId,
  setMeasureExpressionValidation,
} = measureDefSlice.actions;

export const measureDefSliceReducer = measureDefSlice.reducer;
