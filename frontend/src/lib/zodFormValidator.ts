import { z } from "zod"

export function zodValidator<TSchema extends z.ZodTypeAny>(
  schema: TSchema,
) {
  return ({ value }: { value: z.infer<TSchema> }) => {
    const result = schema.safeParse(value)

    if (result.success) return undefined

    const tree = z.treeifyError(result.error)

    return Object.fromEntries(
      Object.entries('properties' in tree ? tree.properties ?? {} : {}).map(([key, val]) => [
        key,
        (val as any)?.errors,
      ]),
    )
  }
}
