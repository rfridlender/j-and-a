import { type PersonMetadata, personMetadataColumns } from "./personMetadata"
import { type Log, logColumns } from "./log"

import type { ColumnDef } from "@tanstack/vue-table"
import { Scroll, User } from "lucide-vue-next"
import type { Component } from "vue"

export type ModelTypes = {
    Log: Log
    PersonMetadata: PersonMetadata
}

type ModelDefinition = {
    icon: Component
    columns: ColumnDef<Log | PersonMetadata>[]
}

export const modelDefinitions: Record<string, ModelDefinition> = {
    log: {
        icon: Scroll,
        columns: logColumns,
    },
    "person-metadata": {
        icon: User,
        columns: personMetadataColumns,
    },
}
