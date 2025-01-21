<script setup lang="ts" generic="TData, TValue">
import { DataTablePagination, DataTableViewOptions } from "@/components/ui/data-table"
import { Input } from "@/components/ui/input"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { valueUpdater } from "@/lib/utils"

import type { ColumnDef, ColumnFiltersState, ExpandedState, SortingState, VisibilityState } from "@tanstack/vue-table"
import {
    FlexRender,
    getCoreRowModel,
    getExpandedRowModel,
    getFilteredRowModel,
    getPaginationRowModel,
    getSortedRowModel,
    useVueTable,
} from "@tanstack/vue-table"
import { ref } from "vue"

const props = defineProps<{
    columns: ColumnDef<TData, TValue>[]
    data: TData[]
}>()

const columnFilters = ref<ColumnFiltersState>([])
const columnVisibility = ref<VisibilityState>({})
const expanded = ref<ExpandedState>({})
const sorting = ref<SortingState>([])
const rowSelection = ref({})

const table = useVueTable({
    get data() {
        return props.data
    },
    get columns() {
        return props.columns
    },
    getCoreRowModel: getCoreRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    onColumnFiltersChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnFilters),
    onColumnVisibilityChange: (updaterOrValue) => valueUpdater(updaterOrValue, columnVisibility),
    onExpandedChange: (updaterOrValue) => valueUpdater(updaterOrValue, expanded),
    onRowSelectionChange: (updaterOrValue) => valueUpdater(updaterOrValue, rowSelection),
    onSortingChange: (updaterOrValue) => valueUpdater(updaterOrValue, sorting),
    state: {
        get sorting() {
            return sorting.value
        },
        get columnFilters() {
            return columnFilters.value
        },
        get columnVisibility() {
            return columnVisibility.value
        },
        get expanded() {
            return expanded.value
        },
        get rowSelection() {
            return rowSelection.value
        },
    },
})
</script>

<template>
    <div>
        <div class="flex items-center py-4">
            <Input
                class="max-w-sm"
                placeholder="Filter last names..."
                :model-value="table.getColumn('familyName')?.getFilterValue() as string"
                @update:model-value="table.getColumn('familyName')?.setFilterValue($event)"
            />

            <DataTableViewOptions :table="table" />
        </div>

        <div class="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
                        <TableHead v-for="header in headerGroup.headers" :key="header.id">
                            <FlexRender
                                v-if="!header.isPlaceholder"
                                :render="header.column.columnDef.header"
                                :props="header.getContext()"
                            />
                        </TableHead>
                    </TableRow>
                </TableHeader>

                <TableBody>
                    <template v-if="table.getRowModel().rows?.length">
                        <template v-for="row in table.getRowModel().rows" :key="row.id">
                            <TableRow :data-state="row.getIsSelected() ? 'selected' : undefined">
                                <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                                    <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
                                </TableCell>
                            </TableRow>

                            <TableRow v-if="row.getIsExpanded()">
                                <TableCell :colspan="row.getAllCells().length">
                                    {{ JSON.stringify(row.original) }}
                                </TableCell>
                            </TableRow>
                        </template>
                    </template>

                    <template v-else>
                        <TableRow>
                            <TableCell :colspan="columns.length" class="h-24 text-center"> No results. </TableCell>
                        </TableRow>
                    </template>
                </TableBody>
            </Table>
        </div>

        <DataTablePagination :table="table" />
    </div>
</template>
