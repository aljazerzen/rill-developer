import { DataModelerStateService } from "./data-modeler-state-service/DataModelerStateService";
import type { DataModelerService } from "./data-modeler-service/DataModelerService";
import { DataModelerSocketService } from "./socket/DataModelerSocketService";
import type { RootConfig } from "./config/RootConfig";
import { PersistentTableEntityService } from "./data-modeler-state-service/entity-state-service/PersistentTableEntityService";
import { DerivedTableEntityService } from "./data-modeler-state-service/entity-state-service/DerivedTableEntityService";
import { PersistentModelEntityService } from "./data-modeler-state-service/entity-state-service/PersistentModelEntityService";
import { DerivedModelEntityService } from "./data-modeler-state-service/entity-state-service/DerivedModelEntityService";
import { ApplicationStateService } from "./data-modeler-state-service/entity-state-service/ApplicationEntityService";
import type { MetricsService } from "./metrics-service/MetricsService";
import { MetricsSocketService } from "./socket/MetricsSocketService";
import { MetricsDefinitionStateService } from "./data-modeler-state-service/entity-state-service/MetricsDefinitionEntityService";
import { MeasureDefinitionStateService } from "./data-modeler-state-service/entity-state-service/MeasureDefinitionStateService";
import { DimensionDefinitionStateService } from "./data-modeler-state-service/entity-state-service/DimensionDefinitionStateService";

export function dataModelerStateServiceClientFactory() {
  return new DataModelerStateService(
    [],
    [
      PersistentTableEntityService,
      DerivedTableEntityService,
      PersistentModelEntityService,
      DerivedModelEntityService,
      ApplicationStateService,
      MetricsDefinitionStateService,
      MeasureDefinitionStateService,
      DimensionDefinitionStateService,
    ].map((EntityStateService) => new EntityStateService())
  );
}

export function clientFactory(config: RootConfig): {
  dataModelerStateService: DataModelerStateService;
  metricsService: MetricsService;
  dataModelerService: DataModelerService;
} {
  const dataModelerStateService = dataModelerStateServiceClientFactory();
  const metricsService = new MetricsSocketService(config);
  const dataModelerService = new DataModelerSocketService(
    dataModelerStateService,
    metricsService,
    config.server
  );

  return { dataModelerStateService, metricsService, dataModelerService };
}
