<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    getRuntimeServiceListFilesQueryKey,
    useRuntimeServiceGetCatalogEntry,
    useRuntimeServiceRenameFileAndReconcile,
  } from "@rilldata/web-common/runtime-client";
  import { EntityType } from "@rilldata/web-local/common/data-modeler-state-service/entity-state-service/EntityStateService";
  import { createForm } from "svelte-forms-lib";
  import * as yup from "yup";
  import {
    dataModelerService,
    runtimeStore,
  } from "../../application-state-stores/application-store";
  import { updateMetricsDefsWrapperApi } from "../../redux-store/metrics-definition/metrics-definition-apis";
  import { store } from "../../redux-store/store-root";
  import { queryClient } from "../../svelte-query/globalQueryClient";
  import Input from "../forms/Input.svelte";
  import SubmissionError from "../forms/SubmissionError.svelte";
  import { Dialog } from "../modal/index";
  import notifications from "../notifications";
  // TODO: get rid of entity id once model and metrics are moved over
  export let entityId = null;

  export let closeModal: () => void;
  export let entityType: EntityType;
  export let currentAssetName: string;

  let error: string;

  $: runtimeInstanceId = $runtimeStore.instanceId;
  $: getCatalog = useRuntimeServiceGetCatalogEntry(
    runtimeInstanceId,
    currentAssetName
  );

  const renameSource = useRuntimeServiceRenameFileAndReconcile();

  const { form, errors, handleSubmit } = createForm({
    initialValues: {
      newName: currentAssetName,
    },
    validationSchema: yup.object({
      newName: yup
        .string()
        .matches(
          /^[a-zA-Z_][a-zA-Z0-9_]*$/,
          "Name must start with a letter or underscore and contain only letters, numbers, and underscores"
        )
        .required("Enter a name!")
        .notOneOf([currentAssetName], `That's the current name!`),
    }),
    onSubmit: (values) => {
      // TODO: remove this branching logic once we have a unified backend for all entities
      switch (entityType) {
        case EntityType.Table: {
          // CHECK: Is this `updateTableName` API call necessary?
          dataModelerService.dispatch("updateTableName", [
            entityId,
            values.newName,
          ]);
          $renameSource.mutate(
            {
              data: {
                instanceId: runtimeInstanceId,
                fromPath: `sources/${currentAssetName}.yaml`,
                toPath: `sources/${values.newName}.yaml`,
              },
            },
            {
              onSuccess: () => {
                closeModal();
                goto(`/source/${values.newName}`, { replaceState: true });
                notifications.send({
                  message: `renamed ${entityLabel} ${currentAssetName} to ${values.newName}`,
                });
                return queryClient.invalidateQueries(
                  getRuntimeServiceListFilesQueryKey($runtimeStore.instanceId)
                );
              },
              onError: (err) => {
                error = err.response.data.message;
                // reset the new table name
                dataModelerService.dispatch("updateTableName", [entityId, ""]);
              },
            }
          );
          break;
        }
        case EntityType.Model: {
          dataModelerService
            .dispatch("updateModelName", [entityId, values.newName])
            .then((response) => {
              if (response.status === 0) {
                notifications.send({ message: response.messages[0].message });
                closeModal();
              } else if (response.status === 1) {
                error = response.messages[0].message;
              }
            });
          break;
        }
        case EntityType.MetricsDefinition: {
          store.dispatch(
            updateMetricsDefsWrapperApi({
              id: entityId,
              changes: { metricDefLabel: values.newName },
            })
          );
          closeModal();
          notifications.send({
            message: `dashboard renamed to ${values.newName}`,
          });
          break;
        }
        default:
          throw new Error(
            "entityType must be either 'Table', 'Model', or 'MetricsDefinition'"
          );
      }
    },
  });

  function getLabel(entityType: EntityType) {
    switch (entityType) {
      case EntityType.Table:
        return "source";
      case EntityType.Model:
        return "model";
      case EntityType.MetricsDefinition:
        return "dashboard";
      default:
        throw new Error(
          "entityType must be either 'Table', 'Model', or 'MetricsDefinition'"
        );
    }
  }

  $: entityLabel = getLabel(entityType);
</script>

<Dialog
  compact
  disabled={$form["newName"] === ""}
  on:cancel={closeModal}
  on:click-outside={closeModal}
  on:primary-action={handleSubmit}
  showCancel
  size="sm"
>
  <svelte:fragment slot="title">Rename</svelte:fragment>
  <div slot="body">
    {#if error}
      <SubmissionError message={error} />
    {/if}
    <form autocomplete="off" on:submit|preventDefault={handleSubmit}>
      <div class="py-2">
        <Input
          bind:value={$form["newName"]}
          claimFocusOnMount
          error={$errors["newName"]}
          id="{entityLabel}-name"
          label="{entityLabel} name"
        />
      </div>
    </form>
  </div>
  <svelte:fragment slot="primary-action-body">Change Name</svelte:fragment>
</Dialog>
