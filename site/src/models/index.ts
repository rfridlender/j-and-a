import type { ColumnDef } from "@tanstack/vue-table"
import { PersonMetadataColumns } from "./PersonMetadata"

export { type PersonMetadata } from "./PersonMetadata"
export const columns: Record<string, never> = {
    PersonMetadata: PersonMetadataColumns,
}
