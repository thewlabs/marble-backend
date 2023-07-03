import * as yup from "yup";
import { adaptDtoWithYup } from "@/infra/adaptDtoWithYup";
import { type Scenario } from "./Scenario";

const ScenarioSchema = yup.object({
  id: yup.string().required(),
  name: yup.string().required(),
  description: yup.string().required(),
  triggerObjectType: yup.string().required(),
  createdAt: yup.date().required(),
  liveVersionId: yup.string().defined().nullable(),
});

export type ScenarioDto = yup.InferType<typeof ScenarioSchema>;

export function adaptScenario(dto: ScenarioDto): Scenario {
  return {
    scenarioId: dto.id,
    name: dto.name,
    description: dto.description,
    triggerObjectType: dto.triggerObjectType,
    createdAt: dto.createdAt,
    liveVersionId: dto.liveVersionId,
  };
}

export function adaptScenariosApiResult(json: unknown): Scenario[] {
  const dtos = adaptDtoWithYup(json, yup.array().required().of(ScenarioSchema));
  return dtos.map(adaptScenario);
}

export function adaptSingleScenarioApiResult(json: unknown): Scenario {
  const dto = adaptDtoWithYup(
    json,
    yup.object({
      scenario: ScenarioSchema,
    })
  );
  return adaptScenario(dto.scenario);
}
