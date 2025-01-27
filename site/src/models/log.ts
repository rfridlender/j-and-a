import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { DataTableDropdown } from "@/components/ui/data-table"

import type { ColumnDef } from "@tanstack/vue-table"
import { ArrowUpDown } from "lucide-vue-next"
import { h } from "vue"

export type Log = {
    givenName: string
    familyName: string
    personId: string
    createdAt: string
    createdBy: string
    deletedAt: string
    deletedBy: string
}

export const logColumns: ColumnDef<Log>[] = [
    {
        id: "select",
        header: ({ table }) =>
            h(Checkbox, {
                checked: table.getIsAllPageRowsSelected(),
                "onUpdate:checked": (value: boolean) => table.toggleAllPageRowsSelected(!!value),
                ariaLabel: "Select all",
            }),
        cell: ({ row }) =>
            h(Checkbox, {
                checked: row.getIsSelected(),
                "onUpdate:checked": (value: boolean) => row.toggleSelected(!!value),
                ariaLabel: "Select row",
            }),
        enableSorting: false,
        enableHiding: false,
    },
    {
        accessorKey: "givenName",
        header: () => h("div", "First Name"),
        cell: ({ row }) => h("div", row.getValue("givenName")),
    },
    {
        accessorKey: "familyName",
        header: ({ column }) => {
            return h(
                Button,
                {
                    variant: "ghost",
                    onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
                },
                () => ["Last Name", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
            )
        },
        cell: ({ row }) => h("div", row.getValue("familyName")),
    },
    {
        id: "actions",
        enableHiding: false,
        cell: ({ row }) =>
            h("div", { class: "relative" }, h(DataTableDropdown, { id: row.original.personId, onExpand: row.toggleExpanded })),
    },
]
