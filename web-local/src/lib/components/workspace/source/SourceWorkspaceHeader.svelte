<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    getRuntimeServiceGetCatalogEntryQueryKey,
    getRuntimeServiceListFilesQueryKey,
    useRuntimeServiceGetCatalogEntry,
    useRuntimeServicePutFileAndReconcile,
    useRuntimeServiceRenameFileAndReconcile,
    useRuntimeServiceTriggerRefresh,
  } from "@rilldata/web-common/runtime-client";
  import { refreshSource } from "@rilldata/web-local/lib/components/navigation/sources/refreshSource";
  import { queryClient } from "@rilldata/web-local/lib/svelte-query/globalQueryClient";
  import { getContext } from "svelte";
  import { fade } from "svelte/transition";
  import {
    dataModelerService,
    runtimeStore,
  } from "../../../application-state-stores/application-store";
  import { overlay } from "../../../application-state-stores/overlay-store";
  import type { PersistentTableStore } from "../../../application-state-stores/table-stores";
  import { IconButton } from "../../button";
  import Import from "../../icons/Import.svelte";
  import RefreshIcon from "../../icons/RefreshIcon.svelte";
  import Source from "../../icons/Source.svelte";
  import notifications from "../../notifications";
  import Tooltip from "../../tooltip/Tooltip.svelte";
  import TooltipContent from "../../tooltip/TooltipContent.svelte";
  import WorkspaceHeader from "../core/WorkspaceHeader.svelte";

  export let id;
  export let name: string;

  const persistentTableStore = getContext(
    "rill:app:persistent-table-store"
  ) as PersistentTableStore;

  $: currentSource = $persistentTableStore?.entities?.find(
    (entity) => entity.id === id || entity.tableName === name
  );

  const renameSource = useRuntimeServiceRenameFileAndReconcile();

  const onChangeCallback = async (e) => {
    if (!e.target.value.match(/^[a-zA-Z_][a-zA-Z0-9_]*$/)) {
      notifications.send({
        message:
          "Source name must start with a letter or underscore and contain only letters, numbers, and underscores",
      });
      e.target.value = currentSource.name; // resets the input
      return;
    }

    dataModelerService.dispatch("updateTableName", [id, e.target.value]);
    $renameSource.mutate(
      {
        data: {
          instanceId: runtimeInstanceId,
          fromPath: `sources/${name}.yaml`,
          toPath: `sources/${e.target.value}.yaml`,
        },
      },
      {
        onSuccess: () => {
          goto(`/source/${e.target.value}`, { replaceState: true });
          return queryClient.invalidateQueries(
            getRuntimeServiceListFilesQueryKey($runtimeStore.instanceId)
          );
        },
        onError: (err) => {
          console.error(err.response.data.message);
          // reset the new table name
          dataModelerService.dispatch("updateTableName", [
            currentSource?.id,
            "",
          ]);
        },
      }
    );
  };

  $: runtimeInstanceId = $runtimeStore.instanceId;
  const refreshSourceMutation = useRuntimeServiceTriggerRefresh();
  const createSource = useRuntimeServicePutFileAndReconcile();

  $: getSource = useRuntimeServiceGetCatalogEntry(
    runtimeInstanceId,
    currentSource?.tableName
  );

  $: connector = $getSource.data?.entry?.source.connector as string;

  const onRefreshClick = async (tableName: string) => {
    try {
      await refreshSource(
        connector,
        tableName,
        $runtimeStore,
        $refreshSourceMutation,
        $createSource
      );
      // invalidate the data preview (async)
      dataModelerService.dispatch("collectTableInfo", [currentSource?.id]);

      // invalidate the "refreshed_on" time
      const queryKey = getRuntimeServiceGetCatalogEntryQueryKey(
        runtimeInstanceId,
        tableName
      );
      await queryClient.invalidateQueries(queryKey);
    } catch (err) {
      // no-op
    }
    overlay.set(null);
  };

  function formatRefreshedOn(refreshedOn: string) {
    const date = new Date(refreshedOn);
    return date.toLocaleString(undefined, {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "numeric",
      minute: "numeric",
    });
  }
</script>

<div class="grid  items-center" style:grid-template-columns="auto max-content">
  <WorkspaceHeader
    {...{ titleInput: name, onChangeCallback }}
    showStatus={false}
  >
    <svelte:fragment slot="icon">
      <Source />
    </svelte:fragment>
    <svelte:fragment slot="right">
      {#if $refreshSourceMutation.isLoading}
        Refreshing...
      {:else}
        <div class="flex items-center">
          {#if $getSource.isSuccess && $getSource.data?.entry?.refreshedOn}
            <div
              class="ui-copy-muted"
              transition:fade|local={{ duration: 200 }}
            >
              Imported on {formatRefreshedOn(
                $getSource.data?.entry?.refreshedOn
              )}
            </div>
          {/if}
          {#if connector === "file"}
            <Tooltip location="bottom" distance={8}>
              <div style="transformY(-1px)">
                <IconButton
                  on:click={() => onRefreshClick(currentSource.tableName)}
                >
                  <Import size="16px" />
                </IconButton>
              </div>
              <TooltipContent slot="tooltip-content">
                Import local file to refresh source
              </TooltipContent>
            </Tooltip>
          {:else}
            <Tooltip location="bottom" distance={8}>
              <IconButton
                on:click={() => onRefreshClick(currentSource.tableName)}
              >
                <RefreshIcon size="16px" />
              </IconButton>
              <TooltipContent slot="tooltip-content">
                Refresh the source data
              </TooltipContent>
            </Tooltip>
          {/if}
        </div>
      {/if}
    </svelte:fragment>
  </WorkspaceHeader>
</div>
