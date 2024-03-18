import { createFileRoute } from '@tanstack/react-router'
import {Login} from "../pages";

export const Route = createFileRoute("/login")({
    component: LoginRoute,
})

function LoginRoute() {
    return <Login></Login>
}