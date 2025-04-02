import { DataTableColumnHeader, DataTableRowSelectorCell, DataTableStrikableCell } from "@/components/ui/data-table"
import { DataTableRowSelectorColumnHeader } from "@/components/ui/data-table"
import { mediumDateFormatter } from "@/lib/dateFormatters"
import { coreSchema } from "@/models/core"

import { parseAbsoluteToLocal } from "@internationalized/date"
import type { ColumnDef } from "@tanstack/vue-table"
import * as changeCase from "change-case"
import { v7 as uuidv7 } from "uuid"
import { h } from "vue"
import { z } from "zod"

export const schema = z.object({
    givenName: z.string().trim().min(1, "Required").default(""),
    familyName: z.string().trim().min(1, "Required").default(""),
    personId: z.string().uuid().default(uuidv7),
})
const mergedSchema = schema.merge(coreSchema)
export type Type = z.infer<typeof mergedSchema>

export function getInitialValues() {
    return {
        givenName: "",
        familyName: "",
        personId: uuidv7(),
    }
}

export const columns: ColumnDef<Type>[] = [
    {
        id: "select",
        header: ({ table }) => h(DataTableRowSelectorColumnHeader, { table: table }),
        cell: ({ row }) => h(DataTableRowSelectorCell, { row: row }),
        enableColumnFilter: false,
        enableHiding: false,
        enableSorting: false,
    },
    {
        accessorKey: "givenName",
        id: "givenName",
        header: ({ column }) => h(DataTableColumnHeader, { column: column, title: changeCase.capitalCase(column.id) }),
        cell: ({ row }) => h(DataTableStrikableCell, { row: row }, () => row.getValue("givenName")),
    },
    {
        accessorKey: "familyName",
        id: "familyName",
        header: ({ column }) => h(DataTableColumnHeader, { column: column, title: changeCase.capitalCase(column.id) }),
        cell: ({ row }) => h(DataTableStrikableCell, { row: row }, () => row.getValue("familyName")),
    },
    {
        accessorFn: (row) => row.deletedAt || row.createdAt,
        id: "updatedAt",
        header: ({ column }) =>
            h(DataTableColumnHeader, { column: column, title: changeCase.capitalCase(column.id), class: "justify-end" }),
        cell: ({ row }) =>
            h(DataTableStrikableCell, { class: "text-right", row: row }, () =>
                mediumDateFormatter.format(parseAbsoluteToLocal(row.getValue("updatedAt")).toDate())
            ),
        enableColumnFilter: false,
    },
]
