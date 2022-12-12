import type { ParamMatcher } from "@sveltejs/kit";

/** @type {import('@sveltejs/kit').ParamMatcher} */
export const match: ParamMatcher = function (param: string) {
    return /^\d+$/.test(param);
}