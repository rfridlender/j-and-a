<script setup lang="ts">
import {
    DataTableFilterOptions,
    DataTablePagination,
    DataTableRowOptions,
    DataTableViewOptions,
} from "@/components/ui/data-table"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { valueUpdater } from "@/lib/utils"
import type { ModelTypeValues } from "@/models"

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
import { LoaderCircle } from "lucide-vue-next"
import { ref, watch } from "vue"

const props = defineProps<{
    columns: ColumnDef<ModelTypeValues>[]
    data: ModelTypeValues[]
    isDataLoading: boolean
}>()

defineEmits<{
    delete: [originalRow: ModelTypeValues]
    duplicate: [originalRow: ModelTypeValues]
    edit: [originalRow: ModelTypeValues]
    restore: [originalRow: ModelTypeValues]
}>()

const columnFilters = ref<ColumnFiltersState>([])
const columnVisibility = ref<VisibilityState>({})
const expanded = ref<ExpandedState>({})
const rowSelection = ref({})
const sorting = ref<SortingState>([])

const filterByColumnId = ref(props.columns.find(({ enableColumnFilter }) => enableColumnFilter !== false)?.id || "")
watch(filterByColumnId, (_, oldValue) => table.getColumn(oldValue)?.setFilterValue(""))

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
        get sorting() {
            return sorting.value
        },
    },
})
</script>

<template>
    <div class="w-full">
        <div class="flex justify-between items-center pb-4">
            <DataTableFilterOptions v-model="filterByColumnId" :table="table" />
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
                            <TableRow class="group relative" :data-state="row.getIsSelected() ? 'selected' : undefined">
                                <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                                    <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
                                </TableCell>

                                <DataTableRowOptions
                                    :row="row"
                                    @delete="(originalRow) => $emit('delete', originalRow)"
                                    @duplicate="(originalRow) => $emit('duplicate', originalRow)"
                                    @edit="(originalRow) => $emit('edit', originalRow)"
                                    @restore="(originalRow) => $emit('restore', originalRow)"
                                />
                            </TableRow>

                            <TableRow v-if="row.getIsExpanded()">
                                <TableCell :colspan="row.getAllCells().length">
                                    {{ JSON.stringify(row.original) }}
                                </TableCell>
                            </TableRow>
                        </template>
                    </template>

                    <template v-else-if="isDataLoading">
                        <TableRow>
                            <TableCell :colspan="columns.length" class="h-24">
                                <div class="flex flex-col justify-center items-center">
                                    <LoaderCircle class="animate-spin" />
                                    Fetching results...
                                </div>
                            </TableCell>
                        </TableRow>
                    </template>

                    <template v-else>
                        <TableRow>
                            <TableCell :colspan="columns.length" class="h-24 text-center">No results.</TableCell>
                        </TableRow>
                    </template>
                </TableBody>
            </Table>
        </div>

        <DataTablePagination :table="table" />
    </div>
</template>
