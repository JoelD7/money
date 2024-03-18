import { createFileRoute } from '@tanstack/react-router'
import { SignUp} from "../pages";

export const Route = createFileRoute("/signup")({
    component: SignUpRoute,
})

function SignUpRoute() {
    return <SignUp></SignUp>
}