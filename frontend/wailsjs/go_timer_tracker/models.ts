// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT


// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as time$0 from "../time/models.js";

export enum ColorTheme {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    Light = "light",
    Dark = "dark",
};

export interface StopwatchData {
    "IsActive": boolean;
    "Hue": number;
    "Name": string;
    "LastUpdated": time$0.Time;
    "CreatedAt": time$0.Time;

    /**
     * Time periods at which stopwatch was ticking, where each period is a 2-element array
     * with first element being the start of the period, and second being the end of it.
     */
    "ActivePeriods": time$0.Time[][] | null;
}
