import {createFileRoute, redirect} from "@tanstack/react-router";
import {Home} from "../pages";
import {store} from "../store";
import {z} from "zod";

function isAuth() {
    return store.getState().authReducer.isAuthenticated;
}

const expensesSearchSchema = z.object({
    categories: z.string().optional(),
    pageSize: z.number().optional(),
    startKey: z.string().optional(),
    period: z.string().optional(),
    sortBy: z.enum(["created_date", "name", "amount"]).optional(),
    sortOrder: z.enum(["asc", "desc"]).optional(),
});

export const Route = createFileRoute("/")({
    beforeLoad: async ({location}) => {
        if (!isAuth()) {
            throw redirect({
                to: "/login",
                search: {
                    redirect: location.href,
                },
            });
        }
    },
    validateSearch: (search) => {
        let searchParams: z.TypeOf<typeof expensesSearchSchema>

        try {
            searchParams = expensesSearchSchema.parse(search)
        } catch (e) {
            console.error("Search params parsing failed:", e)
            throw new Error("Search params parsing failed")
        }

        return searchParams
    },
    component: Index,
});


function Index() {
    return <Home/>;
}
