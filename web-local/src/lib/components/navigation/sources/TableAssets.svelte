<script lang="ts">
  import { page } from "$app/stores";
  import { EntityType } from "@rilldata/web-local/common/data-modeler-state-service/entity-state-service/EntityStateService";
  import { LIST_SLIDE_DURATION } from "@rilldata/web-local/lib/application-config";
  import { useSourceNames } from "@rilldata/web-local/lib/svelte-query/sources";
  import { getContext } from "svelte";
  import { flip } from "svelte/animate";
  import { slide } from "svelte/transition";
  import { runtimeStore } from "../../../application-state-stores/application-store";
  import type { PersistentModelStore } from "../../../application-state-stores/model-stores";
  import type {
    DerivedTableStore,
    PersistentTableStore,
  } from "../../../application-state-stores/table-stores";
  import { createModelForSource } from "../../../redux-store/source/source-apis";
  import ColumnProfile from "../../column-profile/ColumnProfile.svelte";
  import Source from "../../icons/Source.svelte";
  import NavigationEntry from "../NavigationEntry.svelte";
  import NavigationHeader from "../NavigationHeader.svelte";
  import RenameAssetModal from "../RenameAssetModal.svelte";
  import AddSourceModal from "./AddSourceModal.svelte";
  import SourceMenuItems from "./SourceMenuItems.svelte";
  import SourceTooltip from "./SourceTooltip.svelte";

  $: sourceNames = useSourceNames($runtimeStore.instanceId);

  const persistentTableStore = getContext(
    "rill:app:persistent-table-store"
  ) as PersistentTableStore;

  const derivedTableStore = getContext(
    "rill:app:derived-table-store"
  ) as DerivedTableStore;

  const persistentModelStore = getContext(
    "rill:app:persistent-model-store"
  ) as PersistentModelStore;

  let showTables = true;

  let showAddSourceModal = false;

  const openShowAddSourceModal = () => {
    showAddSourceModal = true;
  };

  const queryHandler = async (tableName: string) => {
    const asynchronous = true;
    await createModelForSource(
      $persistentModelStore.entities,
      tableName,
      asynchronous
    );
  };

  let showRenameTableModal = false;
  let renameTableID = null;
  let renameTableName = null;

  const openRenameTableModal = (tableID: string, tableName: string) => {
    showRenameTableModal = true;
    renameTableID = tableID;
    renameTableName = tableName;
  };
</script>

<NavigationHeader
  bind:show={showTables}
  tooltipText="add a new data source"
  on:add={openShowAddSourceModal}
  contextButtonID={"add-table"}
>
  <Source size="16px" /> Sources
</NavigationHeader>

{#if showTables}
  <div class="pb-6" transition:slide|local={{ duration: LIST_SLIDE_DURATION }}>
    {#if $sourceNames?.data && $persistentTableStore?.entities && $derivedTableStore?.entities}
      <!-- TODO: fix the object property access back to t.id from t["id"] once svelte fixes it -->
      {#each $sourceNames.data as sourceName (sourceName)}
        {@const persistentTable = $persistentTableStore.entities.find(
          (t) => t["tableName"] == sourceName
        )}
        {@const derivedTable = $derivedTableStore.entities.find(
          (t) => t["id"] === persistentTable?.id
        )}
        <div
          animate:flip={{ duration: 200 }}
          out:slide={{ duration: LIST_SLIDE_DURATION }}
        >
          <NavigationEntry
            href={`/source/${sourceName}`}
            open={$page.url.pathname === `/source/${sourceName}`}
            on:command-click={() => queryHandler(sourceName)}
            name={sourceName}
          >
            <svelte:fragment slot="more">
              <div transition:slide|local={{ duration: LIST_SLIDE_DURATION }}>
                <ColumnProfile
                  indentLevel={1}
                  cardinality={derivedTable?.cardinality ?? 0}
                  profile={derivedTable?.profile ?? []}
                  head={derivedTable?.preview ?? []}
                  entityId={persistentTable.id}
                />
              </div>
            </svelte:fragment>

            <svelte:fragment slot="tooltip-content">
              <SourceTooltip {sourceName} />
            </svelte:fragment>

            <svelte:fragment slot="menu-items" let:toggleMenu>
              <SourceMenuItems
                {sourceName}
                {toggleMenu}
                on:rename-asset={() => {
                  openRenameTableModal(persistentTable?.id, sourceName);
                }}
              />
            </svelte:fragment>
          </NavigationEntry>
        </div>
      {/each}
    {/if}
  </div>
{/if}

{#if showAddSourceModal}
  <AddSourceModal
    on:close={() => {
      showAddSourceModal = false;
    }}
  />
{/if}
{#if showRenameTableModal}
  <RenameAssetModal
    entityType={EntityType.Table}
    closeModal={() => (showRenameTableModal = false)}
    entityId={renameTableID}
    currentAssetName={renameTableName}
  />
{/if}
