import {
    columns as personMetadataColumns,
    getInitialValues as getPersonMetadataInitialValues,
    schema as personMetadataSchema,
    type Type as PersonMetadata,
} from "@/models/personMetadata"

import type { ColumnDef } from "@tanstack/vue-table"
import { User } from "lucide-vue-next"
import type { Component } from "vue"
import { z } from "zod"

export type ModelTypes = {
    PersonMetadata: PersonMetadata
}

export type ModelType = keyof ModelTypes

type ModelSchema<T> = z.ZodObject<z.ZodRawShape, z.UnknownKeysParam, z.ZodTypeAny, T>

export const definitions = {
    "person-metadata": {
        columns: personMetadataColumns,
        schema: personMetadataSchema as unknown as ModelSchema<PersonMetadata>,
        icon: User,
        partitionType: "Person",
        partitionIdKey: "personId",
        sortType: "PersonMetadata",
        getInitialValues: getPersonMetadataInitialValues,
    },
} satisfies Record<
    string,
    {
        columns: ColumnDef<ModelTypes[ModelType]>[]
        schema: ModelSchema<ModelTypes[ModelType]>
        icon: Component
        partitionType: string
        partitionIdKey: string
        sortType: ModelType
        sortIdKey?: string
        getInitialValues: () => ModelTypes[ModelType]
    }
>
