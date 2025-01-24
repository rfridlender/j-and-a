<script setup lang="ts">
import icon from "@/assets/icon.png"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Separator } from "@/components/ui/separator"
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarInset,
    SidebarMenu,
    SidebarMenuAction,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarMenuSub,
    SidebarMenuSubButton,
    SidebarMenuSubItem,
    SidebarProvider,
    SidebarRail,
    SidebarTrigger,
} from "@/components/ui/sidebar"
import { MODELS, privateRoutes } from "@/router"
import { useUserAttributes } from "@/stores/userAttributes"

import { signOut } from "aws-amplify/auth"
import {
    BadgeCheck,
    ChevronRight,
    ChevronsUpDown,
    Folder,
    Forward,
    LogOut,
    MoreHorizontal,
    SquareTerminal,
    Trash2,
} from "lucide-vue-next"
import { RouterLink, useRouter } from "vue-router"

const router = useRouter()
const { userAttributes } = useUserAttributes()

function titleize(notTitle: string) {
    return notTitle
        .replace(/([a-z])([A-Z])/g, "$1 $2")
        .replace(/([A-Z])([A-Z][a-z])/g, "$1 $2")
        .replace(/^./, (firstCharacter) => firstCharacter.toUpperCase())
        .replace(/\s+/g, " ")
        .trim()
}

const data = {
    user: {
        name: `${userAttributes?.given_name} ${userAttributes?.family_name}`,
        email: userAttributes?.email,
        avatar: `https://api.dicebear.com/9.x/adventurer/svg?seed=${encodeURIComponent(userAttributes?.given_name || "")}`,
    },
    navMain: [
        {
            title: "Placeholder",
            icon: SquareTerminal,
            isActive: true,
            routes: privateRoutes,
        },
    ],
}

async function onSignOut() {
    try {
        await signOut()
        router.replace("/sign-in")
    } catch (error) {
        console.error(error)

        if (error instanceof Error) {
            router.replace({
                path: "/sign-in",
                query: { "error-message": error.message },
            })
        }
    }
}
</script>

<template>
    <SidebarProvider>
        <Sidebar collapsible="icon">
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <SidebarMenuButton
                            size="lg"
                            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground cursor-auto"
                        >
                            <div
                                class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
                            >
                                <img class="size-4" :src="icon" alt="Placeholder" />
                            </div>

                            <div class="grid flex-1 text-left text-sm leading-tight">
                                <span class="truncate font-semibold">Placeholder</span>
                                <span class="truncate text-xs">A placeholder company</span>
                            </div>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>

            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupLabel>Applications</SidebarGroupLabel>

                    <SidebarMenu>
                        <Collapsible as-child :default-open="true" class="group/collapsible">
                            <SidebarMenuItem>
                                <CollapsibleTrigger as-child>
                                    <SidebarMenuButton tooltip="Placeholder">
                                        <component :is="SquareTerminal" />

                                        <span>Placeholder</span>

                                        <ChevronRight
                                            class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"
                                        />
                                    </SidebarMenuButton>
                                </CollapsibleTrigger>

                                <CollapsibleContent>
                                    <SidebarMenuSub>
                                        <SidebarMenuSubItem
                                            v-for="privateRoute in privateRoutes.filter(({ path }) => path !== '/:modelType')"
                                            :key="privateRoute.path"
                                        >
                                            <SidebarMenuSubButton as-child>
                                                <RouterLink :to="privateRoute.path">
                                                    <span>{{ privateRoute.name }}</span>
                                                </RouterLink>
                                            </SidebarMenuSubButton>
                                        </SidebarMenuSubItem>
                                    </SidebarMenuSub>
                                </CollapsibleContent>
                            </SidebarMenuItem>
                        </Collapsible>
                    </SidebarMenu>
                </SidebarGroup>

                <SidebarGroup class="group-data-[collapsible=icon]:hidden">
                    <SidebarGroupLabel>Models</SidebarGroupLabel>

                    <SidebarMenu>
                        <SidebarMenuItem v-for="model in MODELS" :key="model.modelType">
                            <SidebarMenuButton as-child>
                                <RouterLink :to="`/${model.modelType}`">
                                    <component :is="model.icon" />
                                    <span>{{ titleize(model.modelType) }}</span>
                                </RouterLink>
                            </SidebarMenuButton>

                            <DropdownMenu>
                                <DropdownMenuTrigger as-child>
                                    <SidebarMenuAction show-on-hover>
                                        <MoreHorizontal />
                                        <span class="sr-only">More</span>
                                    </SidebarMenuAction>
                                </DropdownMenuTrigger>

                                <DropdownMenuContent class="w-48 rounded-lg" side="bottom" align="end">
                                    <DropdownMenuItem>
                                        <Folder class="text-muted-foreground" />
                                        <span>View Project</span>
                                    </DropdownMenuItem>

                                    <DropdownMenuItem>
                                        <Forward class="text-muted-foreground" />
                                        <span>Share Project</span>
                                    </DropdownMenuItem>

                                    <DropdownMenuSeparator />

                                    <DropdownMenuItem>
                                        <Trash2 class="text-muted-foreground" />
                                        <span>Delete Project</span>
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroup>
            </SidebarContent>

            <SidebarFooter>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <DropdownMenu>
                            <DropdownMenuTrigger as-child>
                                <SidebarMenuButton
                                    size="lg"
                                    class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                                >
                                    <Avatar class="h-8 w-8 rounded-lg">
                                        <AvatarImage :src="data.user.avatar" :alt="data.user.name" />

                                        <AvatarFallback class="rounded-lg">
                                            {{
                                                `${userAttributes?.given_name?.slice(0, 1)}${userAttributes?.family_name?.slice(0, 1)}`
                                            }}
                                        </AvatarFallback>
                                    </Avatar>

                                    <div class="grid flex-1 text-left text-sm leading-tight">
                                        <span class="truncate font-semibold">{{ data.user.name }}</span>
                                        <span class="truncate text-xs">{{ data.user.email }}</span>
                                    </div>

                                    <ChevronsUpDown class="ml-auto size-4" />
                                </SidebarMenuButton>
                            </DropdownMenuTrigger>

                            <DropdownMenuContent
                                class="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                                side="bottom"
                                align="end"
                                :side-offset="4"
                            >
                                <DropdownMenuLabel class="p-0 font-normal">
                                    <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                                        <Avatar class="h-8 w-8 rounded-lg">
                                            <AvatarImage :src="data.user.avatar" :alt="data.user.name" />

                                            <AvatarFallback class="rounded-lg">
                                                {{
                                                    `${userAttributes?.given_name?.slice(0, 1)}${userAttributes?.family_name?.slice(0, 1)}`
                                                }}
                                            </AvatarFallback>
                                        </Avatar>

                                        <div class="grid flex-1 text-left text-sm leading-tight">
                                            <span class="truncate font-semibold">{{ data.user.name }}</span>
                                            <span class="truncate text-xs">{{ data.user.email }}</span>
                                        </div>
                                    </div>
                                </DropdownMenuLabel>

                                <DropdownMenuSeparator />

                                <DropdownMenuGroup>
                                    <DropdownMenuItem>
                                        <BadgeCheck />
                                        Account
                                    </DropdownMenuItem>
                                </DropdownMenuGroup>

                                <DropdownMenuSeparator />

                                <DropdownMenuItem @click="onSignOut">
                                    <LogOut />
                                    Sign Out
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarFooter>

            <SidebarRail />
        </Sidebar>

        <SidebarInset>
            <header
                class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12"
            >
                <div class="flex items-center gap-2 px-4">
                    <SidebarTrigger class="-ml-1" />

                    <Separator orientation="vertical" class="mr-2 h-4" />

                    <Breadcrumb>
                        <BreadcrumbList>
                            <BreadcrumbItem class="hidden md:block">
                                <BreadcrumbLink href="#">Parent</BreadcrumbLink>
                            </BreadcrumbItem>

                            <BreadcrumbSeparator class="hidden md:block" />

                            <BreadcrumbItem>
                                <BreadcrumbPage>Child</BreadcrumbPage>
                            </BreadcrumbItem>
                        </BreadcrumbList>
                    </Breadcrumb>
                </div>
            </header>

            <RouterView />
        </SidebarInset>
    </SidebarProvider>
</template>
