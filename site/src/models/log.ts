import {
    DataTableColumnHeader,
    DataTableRowSelectorCell,
    DataTableRowSelectorColumnHeader,
    DataTableStrikableCell,
} from "@/components/ui/data-table"
import { mediumDateFormatter } from "@/lib/dateFormatters"
import { coreSchema } from "@/models/core"

import { parseAbsoluteToLocal } from "@internationalized/date"
import type { ColumnDef } from "@tanstack/vue-table"
import * as changeCase from "change-case"
import { h } from "vue"
import { z } from "zod"

export const logSchema = z.object({
    personId: z.string().uuid(),
    hours: z.number().min(0),
    jobId: z.string().uuid(),
    logId: z.string().uuid(),
})
const mergedLogSchema = logSchema.merge(coreSchema)
export type Log = z.infer<typeof mergedLogSchema>

export const logColumns: ColumnDef<Log>[] = [
    {
        id: "select",
        header: ({ table }) => h(DataTableRowSelectorColumnHeader, { table: table }),
        cell: ({ row }) => h(DataTableRowSelectorCell, { row: row }),
        enableColumnFilter: false,
        enableHiding: false,
        enableSorting: false,
    },
    {
        accessorKey: "personId",
        id: "personId",
        header: ({ column }) => h(DataTableColumnHeader, { column: column, title: changeCase.capitalCase(column.id) }),
        cell: ({ row }) => h(DataTableStrikableCell, { row: row }, () => row.getValue("personId")),
    },
    {
        accessorKey: "hours",
        id: "hours",
        header: ({ column }) => h(DataTableColumnHeader, { column: column, title: changeCase.capitalCase(column.id) }),
        cell: ({ row }) => h(DataTableStrikableCell, { row: row }, () => row.getValue("hours")),
    },
    {
        accessorKey: "jobId",
        id: "jobId",
        header: ({ column }) => h(DataTableColumnHeader, { column: column, title: changeCase.capitalCase(column.id) }),
        cell: ({ row }) => h(DataTableStrikableCell, { row: row }, () => row.getValue("jobId")),
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
