import PrivateLayout from "@/layouts/PrivateLayout.vue"
import PublicLayout from "@/layouts/PublicLayout.vue"
import { useAuthSession } from "@/stores/authSession"
import { useUserAttributes } from "@/stores/userAttributes"
import SignIn from "@/views/SignIn.vue"

import { fetchAuthSession, fetchUserAttributes } from "aws-amplify/auth"

import { createRouter, createWebHistory } from "vue-router"

const publicRoutes = [
    {
        path: "/sign-in",
        name: "sign-in",
        component: SignIn,
    },
]

const privateRoutes = [
    {
        path: "/dashboard",
        name: "dashboard",
        component: SignIn,
    },
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        { path: "", component: PublicLayout, children: publicRoutes },
        { path: "", component: PrivateLayout, children: privateRoutes },
    ],
})

router.beforeEach(async (to) => {
    try {
        const authSession = await fetchAuthSession()
        useAuthSession().authSession = authSession

        const userAttributes = await fetchUserAttributes()
        useUserAttributes().userAttributes = userAttributes

        if (to.path === "/" || publicRoutes.some((publicRoute) => publicRoute.path === to.path)) {
            return "/dashboard"
        }
    } catch (error) {
        if (!publicRoutes.some((publicRoute) => publicRoute.path === to.path)) {
            return "/sign-in"
        }
    }
})

export default router
