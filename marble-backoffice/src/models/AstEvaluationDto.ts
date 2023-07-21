import * as yup from "yup";
import { adaptDtoWithYup } from "@/infra/adaptDtoWithYup";
import { type ConstantType } from "./AstExpression";
import {
  ConstantOptionalSchema,
  adaptConstantOptional,
} from "./AstExpressionDto";
import { type AstNodeEvaluation } from "./AstEvaluation";
import { MapObjectValues } from "@/MapUtils";

// Yup can't infer typescript type from recursive schema, let's declare it manually
export interface AstNodeEvaluationDto {
  return_value?: ConstantType;
  evaluation_error: string;
  children?: AstNodeEvaluationDto[];
  named_children?: { [key: string]: AstNodeEvaluationDto };
}

export const AstNodeEvaluationSchema = yup.object({
  return_value: ConstantOptionalSchema,
  evaluation_error: yup.string(),
  children: yup
    .array()
    .of(yup.lazy(() => AstNodeEvaluationSchema.default(null)))
    .optional(),
  named_children: yup.lazy((obj) => {
    return yup.object(
      MapObjectValues(obj || {}, () => {
        return AstNodeEvaluationSchema.default(undefined);
      })
    );
  }),
}) as yup.Schema<AstNodeEvaluationDto>; // Can't use lazy schema as array().of argument in TypeScript: https://github.com/jquense/yup/issues/1190

export function adaptNodeEvaluation(
  dto: AstNodeEvaluationDto
): AstNodeEvaluation {
  return {
    returnValue: adaptConstantOptional(dto.return_value),
    evaluationError: dto.evaluation_error || "",
    children: (dto.children || []).map(adaptNodeEvaluation),
    namedChildren: MapObjectValues(
      dto.named_children || {},
      adaptNodeEvaluation
    ),
  };
}

export function adaptDryRunResult(json: unknown): AstNodeEvaluation {
  const dto = adaptDtoWithYup(
    json,
    yup.object({
      evaluation: AstNodeEvaluationSchema,
    })
  );

  return adaptNodeEvaluation(dto.evaluation);
}
