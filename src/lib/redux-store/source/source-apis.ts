import { dataModelerService } from "$lib/application-state-stores/application-store";
import { EntityType } from "$common/data-modeler-state-service/entity-state-service/EntityStateService";
import notificationStore from "$lib/components/notifications";
import { store } from "$lib/redux-store/store-root";
import {
  createMetricsDefsApi,
  generateMeasuresAndDimensionsApi,
} from "$lib/redux-store/metrics-definition/metrics-definition-apis";
import type { DerivedTableEntity } from "$common/data-modeler-state-service/entity-state-service/DerivedTableEntityService";
import { TIMESTAMPS } from "$lib/duckdb-data-types";
import type { PersistentModelEntity } from "$common/data-modeler-state-service/entity-state-service/PersistentModelEntityService";
import { selectMetricsDefinitionMatchingName } from "$lib/redux-store/metrics-definition/metrics-definition-selectors";
import {
  resetQuickStartDashboardOverlay,
  showQuickStartDashboardOverlay,
} from "$lib/application-state-stores/layout-store";

// Source doesn't have a slice as of now.
// This file has simple code that will eventually be moved into async thunks

/**
 * Create a model for the given source by selecting all columns.
 */
export const createModelForSource = async (
  models: Array<PersistentModelEntity>,
  sourceName: string
) => {
  // change the active asset to the new model
  await dataModelerService.dispatch("setActiveAsset", [
    EntityType.Model,
    await createModelFromSourceAndGetId(models, sourceName),
  ]);

  notificationStore.send({
    message: `queried ${sourceName} in workspace`,
  });
};

/**
 * Quick starts a metrics dashboard for a given source.
 * The source should have a timestamp column for this to work.
 */
export const autoCreateMetricsDefinitionForSource = async (
  models: Array<PersistentModelEntity>,
  derivedSources: Array<DerivedTableEntity>,
  id: string,
  sourceName: string
) => {
  try {
    const timestampColumns = derivedSources
      .find((source) => source.id === id)
      .profile?.filter((column) => TIMESTAMPS.has(column.type));
    if (!timestampColumns?.length) return;
    showQuickStartDashboardOverlay(sourceName, timestampColumns[0].name);
    const modelId = await createModelFromSourceAndGetId(models, sourceName);

    await autoCreateMetricsDefinitionForModel(
      sourceName,
      modelId,
      timestampColumns[0].name
    );
  } catch (e) {
    console.error(e);
  }
  resetQuickStartDashboardOverlay();
};

/**
 * Creates a metrics definition for a given model, time dimension and a label.
 * Auto generates measures and dimensions.
 * Focuses the dashboard created.
 */
export const autoCreateMetricsDefinitionForModel = async (
  sourceName: string,
  sourceModelId: string,
  timeDimension: string
) => {
  const metricsLabel = `metrics_${sourceName}`;
  const existingMetrics = selectMetricsDefinitionMatchingName(
    store.getState(),
    metricsLabel
  );

  const { payload: createdMetricsDef } = await store.dispatch(
    createMetricsDefsApi({
      sourceModelId,
      timeDimension,
      metricDefLabel:
        existingMetrics.length === 0
          ? metricsLabel
          : `${metricsLabel}_${existingMetrics.length}`,
    })
  );

  await store.dispatch(generateMeasuresAndDimensionsApi(createdMetricsDef.id));
  await dataModelerService.dispatch("setActiveAsset", [
    EntityType.MetricsExplorer,
    createdMetricsDef.id,
  ]);
};

/**
 * Create a model with all columns from the source
 */
const createModelFromSourceAndGetId = async (
  models: Array<PersistentModelEntity>,
  sourceName: string
): Promise<string> => {
  // check existing models to avoid a name conflict
  const existingNames = models
    .filter((model) => model.name.includes(`query_${sourceName}`))
    .map((model) => model.tableName)
    .sort();
  const nextName =
    existingNames.length === 0
      ? `query_${sourceName}`
      : `query_${sourceName}_${existingNames.length + 1}`;

  const response = await dataModelerService.dispatch("addModel", [
    {
      name: nextName,
      query: `select * from ${sourceName}`,
    },
  ]);
  return (response as unknown as { id: string }).id;
};