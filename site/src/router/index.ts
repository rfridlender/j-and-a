import PrivateLayout from "@/layouts/PrivateLayout.vue"
import PublicLayout from "@/layouts/PublicLayout.vue"
import { useAuthSession } from "@/stores/authSession"
import { useUserAttributes } from "@/stores/userAttributes"
import ConfirmSignInWithEmailCode from "@/views/ConfirmSignInWithEmailCode.vue"
import DashboardView from "@/views/DashboardView.vue"
import ModelView from "@/views/ModelView.vue"
import SignIn from "@/views/SignIn.vue"

import { fetchAuthSession, fetchUserAttributes } from "aws-amplify/auth"
import { ScrollText, User } from "lucide-vue-next"
import {
    createRouter,
    createWebHistory,
    type NavigationGuardNext,
    type RouteLocationNormalized,
    type RouteLocationNormalizedLoaded,
} from "vue-router"

export const MODELS = [
    { modelType: "Log", icon: ScrollText },
    { modelType: "PersonMetadata", icon: User },
]

const publicRoutes = [
    {
        path: "/sign-in",
        name: "Sign-In",
        component: SignIn,
    },
    {
        path: "/confirm-sign-in-with-email-code",
        name: "Confirm Sign-In With Email Code",
        component: ConfirmSignInWithEmailCode,
        beforeEnter: (to: RouteLocationNormalized, from: RouteLocationNormalizedLoaded, next: NavigationGuardNext) =>
            from.path === "/sign-in" && to.query.email ? next() : next(false),
    },
]

export const privateRoutes = [
    {
        path: "/dashboard",
        name: "Dashboard",
        component: DashboardView,
    },
    {
        path: "/:modelType",
        name: "Model",
        component: ModelView,
        beforeEnter: (to: RouteLocationNormalized, _: RouteLocationNormalizedLoaded, next: NavigationGuardNext) =>
            MODELS.some(
                ({ modelType }) =>
                    modelType === (typeof to.params.modelType === "string" ? to.params.modelType : to.params.modelType[0])
            )
                ? next()
                : next(false),
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
