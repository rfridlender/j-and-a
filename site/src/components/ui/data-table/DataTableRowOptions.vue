<script setup lang="ts">
import { Button } from "@/components/ui/button"
import { TableCell } from "@/components/ui/table"
import type { ModelTypeValues } from "@/models"

import type { Row } from "@tanstack/vue-table"
import { ArchiveRestore, Copy, Pencil, Trash2 } from "lucide-vue-next"

const props = defineProps<{
    row: Row<ModelTypeValues>
}>()

defineEmits<{
    delete: [originalRow: ModelTypeValues]
    duplicate: [originalRow: ModelTypeValues]
    edit: [originalRow: ModelTypeValues]
    restore: [originalRow: ModelTypeValues]
}>()
</script>

<template>
    <TableCell
        class="h-[calc(100%-1px)] pl-12 hidden group-data-[state=selected]:hidden group-hover:table-cell absolute bottom-0 right-0 bg-gradient-to-r from-transparent to-[#18181A] to-25% text-muted-foreground"
    >
        <div class="h-full flex justify-center items-center">
            <Button
                v-if="!props.row.original.deletedAt"
                variant="ghost"
                class="hover:bg-accent hover:text-accent-foreground"
                size="icon"
                @click="$emit('edit', props.row.original)"
            >
                <Pencil />
            </Button>

            <Button
                variant="ghost"
                class="hover:bg-accent hover:text-accent-foreground"
                size="icon"
                @click="$emit('duplicate', props.row.original)"
            >
                <Copy />
            </Button>

            <Button
                v-if="props.row.original.deletedAt"
                variant="ghost"
                class="hover:bg-accent hover:text-accent-foreground"
                size="icon"
                @click="$emit('restore', props.row.original)"
            >
                <ArchiveRestore />
            </Button>

            <Button
                v-else
                variant="ghost"
                class="hover:bg-accent hover:text-accent-foreground"
                size="icon"
                @click="$emit('delete', props.row.original)"
            >
                <Trash2 />
            </Button>
        </div>
    </TableCell>
</template>
