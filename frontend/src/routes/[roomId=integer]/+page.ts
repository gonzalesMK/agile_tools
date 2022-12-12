// since there's no dynamic data here, we can prerender

import type { PageLoad } from "../$types";

// it so that it gets served as a static asset in production
export const prerender = true;

export const load: PageLoad = async ({ params }) => {
    // const res = await fetch('/my-server-route.json');
    //   return await res.json();
    return {
        ...params
    }
}
