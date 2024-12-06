<script lang="ts">
    import { onDestroy, onMount, untrack } from "svelte";
    import { AddStopwatch, DeleteStopwatch, GetStopwatch, HasStopwatch, UpdateStopwatchName, UpdateStopwatchTime } from "@wails/go_timer_tracker/app";
    import Trash from "@assets/icons/Trash.svelte";
    import { tweened } from "svelte/motion";
    import Reset from "@assets/icons/Reset.svelte";

    type Props = {
        id: string,
        onLoaded?: () => any,
        onDelete?: () => any
    }
    let { id, onLoaded = () => {}, onDelete: onDelete = () => {} }: Props = $props();

    let time = $state(0);
    let displayTime: {upper?: string, center: string, lower?: string} = $derived.by(() => {
        let date = new Date(time);

        const hour = 1000 * 60 * 60
        const day = hour * 24

        let upper: string | undefined = undefined;
        if (time / day > 1) {
            upper = `${Math.floor(time / day)} day${time / day > 2 ? 's' : ''}`
        }
        
        let center: string;
        let lower: string | undefined = undefined;
        if (time / hour > 1) {
            center = (
                String(date.getUTCHours()).padStart(2, "0") + 
                ':' +
                String(date.getUTCMinutes()).padStart(2, "0")
            );
            lower = String(date.getUTCSeconds()).padStart(2, "0");
        } else {
            center = (
                String(date.getUTCMinutes()).padStart(2, "0") + 
                ':' +
                String(date.getUTCSeconds()).padStart(2, "0")
            );
        }
        
        return {
            upper,
            center,
            lower
        }
    })

    let active = $state(false);
    let hue = $state(Math.random() * 360);
    let name = $state("Name");

    let initialArrowRotation = $state(0);

    let loading = $state(true);

    let destructionProgress = tweened(0);
    let resetProgress = tweened(0);

    let stopwatchWidth: number = $state(0);

    $effect(() => {
        if (!loading) {
            onLoaded();
        }
    })

    onMount(async () => {
        if (await HasStopwatch(id)) {
            let data = await GetStopwatch(id);
            hue = data.Hue;
            name = data.Name;
            active = data.IsActive;
            time = data.TimeAccumulated;
            if (active) {
                // If timer is active, then it kept going since the last time 
                // we updated the timer backend
                time += Date.now() - Date.parse(data.LastUpdated)
            }
            initialArrowRotation = (time % 1000) / 1000 * 360;
            loading = false;
        } else {
            loading = false;
            await AddStopwatch(id, name, hue)
        }
    })

    $effect(() => {
        if (!active) {
            return;
        }

        let lastTimeChange = Date.now();
        let intervalId: number | undefined = undefined
        let timeoutId = setTimeout(() => {
            time += 1000 - time % 1000;

            lastTimeChange = Date.now();
            intervalId = setInterval(() => {
                time += 1000;
                lastTimeChange = Date.now();
            }, 1000);
        }, untrack(() =>1000 - time % 1000));

        return () => {
            clearTimeout(timeoutId);
            if (intervalId !== undefined) {
                clearInterval(intervalId);
            }
            time += Date.now() - lastTimeChange;
        }
    })

    onDestroy(() => {
        UpdateStopwatchTime(id, active, time);
    })

    async function toggleStopwatch() {
        active = !active;
        await UpdateStopwatchTime(id, active, time)
    }

    function startDestruction() {
        destructionProgress.set(1).then(
            () => {
                onDelete();
                DeleteStopwatch(id)
            }
        );
    }

    function cancelDestruction() {
        destructionProgress.set(0);
    }

    function startReset() {
        resetProgress.set(1).then(
            () => {
                time = time % 1000;
                UpdateStopwatchTime(id, active, time);
                cancelReset();
            }
        );
    }

    function cancelReset() {
        resetProgress.set(0);
    }

    function getOptimalNameSize(): number {
        return stopwatchWidth / (name.length);
    }
</script>

<div class="external" style="--hue: {hue}deg">
    <div class="front-layer" bind:clientWidth={stopwatchWidth}>
        <div class="buttons-bar">
            <button onmousedown={startDestruction} onmouseup={cancelDestruction}>
                <Trash style="width: 100%; height: 100%"></Trash>
            </button>
    
            <button onmousedown={startReset} onmouseup={cancelReset}>
                <Reset style="width: 100%; height: 100%"></Reset>
            </button>
        </div>

        <input type="text" bind:value={name} onfocusout={() => UpdateStopwatchName(id, name)}
                style="font-size: clamp(var(--size-6), {getOptimalNameSize()}px, var(--size-8));"
                size="{Math.max(name.length, 4)}" class="name"/>
    </div>

    <div class="progress-bar destruction" style="--progress: {$destructionProgress * 100}%"></div>

    <div class="progress-bar reset" style="--progress: {$resetProgress * 100}%"></div>

    <button class="outlined display" class:active onclick={toggleStopwatch}
            style="--hue: {hue}deg; --initial-rotation: {initialArrowRotation}deg">
        <div class="background">
            <!-- <svg class="stopwatch-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 56 56">
                <path d="M28 54.273c13.055 0 23.906-10.828 23.906-23.906c0-5.46-1.898-10.523-5.062-14.578l1.828-1.851c.539-.54.844-1.266.844-1.899c0-1.312-1.055-2.367-2.415-2.367c-.796 0-1.359.234-1.921.773l-1.782 1.758c-3.586-3.047-8.039-5.086-12.914-5.601V4.305c0-1.407-1.125-2.578-2.508-2.578c-1.382 0-2.53 1.171-2.53 2.578v2.297C13.538 7.89 4.093 18.156 4.093 30.367c0 13.078 10.828 23.906 23.906 23.906m0-3.984C16.937 50.29 8.1 41.43 8.1 30.367c0-11.039 8.813-19.922 19.876-19.922c11.039 0 19.921 8.883 19.945 19.922C47.945 41.43 39.039 50.29 28 50.29m0-16.125a3.75 3.75 0 0 0 3.75-3.773a3.77 3.77 0 0 0-2.133-3.399V15.344c0-.914-.727-1.617-1.64-1.617c-.915 0-1.618.703-1.618 1.617v11.625c-1.289.68-2.156 1.945-2.156 3.422A3.776 3.776 0 0 0 28 34.164" />
            </svg> -->
            <svg class="stopwatch-icon" viewBox="0 0 56 56" version="1.1" id="svg4" xmlns="http://www.w3.org/2000/svg">
                <defs
                    id="defs8" />
                <path class="arrow" class:active
                    d="m 28,34.165 a 3.75,3.75 0 0 0 3.75,-3.773 3.77,3.77 0 0 0 -2.133,-3.399 V 15.344 c 0,-0.914 -0.727,-1.617 -1.64,-1.617 -0.915,0 -1.618,0.703 -1.618,1.617 v 11.625 c -1.289,0.68 -2.156,1.945 -2.156,3.422 A 3.776,3.776 0 0 0 28,34.164"
                    id="path346" />
                <path
                    d="m 28,54.273 c 13.055,0 23.906,-10.828 23.906,-23.906 0,-5.46 -1.898,-10.523 -5.062,-14.578 l 1.828,-1.851 c 0.539,-0.54 0.844,-1.266 0.844,-1.899 0,-1.312 -1.055,-2.367 -2.415,-2.367 -0.796,0 -1.359,0.234 -1.921,0.773 l -1.782,1.758 C 39.812,9.156 35.359,7.117 30.484,6.602 V 4.305 c 0,-1.407 -1.125,-2.578 -2.508,-2.578 -1.382,0 -2.53,1.171 -2.53,2.578 V 6.602 C 13.538,7.89 4.093,18.156 4.093,30.367 c 0,13.078 10.828,23.906 23.906,23.906 m 0,-3.984 C 16.937,50.29 8.1,41.43 8.1,30.367 c 0,-11.039 8.813,-19.922 19.876,-19.922 11.039,0 19.921,8.883 19.945,19.922 C 47.945,41.43 39.039,50.29 28,50.29"
                    id="path2" />
            </svg>
    
        </div>

        <svg class="time" viewBox="0 0 20 15">
            <text style="font-size: 8.5px" text-anchor="middle" x="10" y="10">
                {displayTime.center}
                <!-- 06:28 -->
            </text>
            {#if displayTime.upper !== undefined}
                <text style="font-size: 4px" text-anchor="middle" x="10" y="4">
                    {displayTime.upper}
                    <!-- 23 days -->
                </text>
            {/if}
            {#if displayTime.lower !== undefined}    
                <text style="font-size: 4px" text-anchor="middle" x="10" y="14">
                    {displayTime.lower}
                    <!-- 47 -->
                </text>
            {/if}
        </svg>
    </button>
</div>

<style>
    .front-layer {
        position: absolute;
        width: 100%;
        height: 100%;
        display: flex;
        flex-flow: column;
        justify-content: space-between;
    }

    .name {
        background: transparent;
        max-width: 100%;
        margin-left: auto;
        margin-right: auto;
        text-align: center;
        z-index: 2;
        text-overflow: ellipsis;
        line-height: 0;
        margin-bottom: var(--size-1);
    }

    :root.dark .name {
        color: hsl(var(--hue), 100%, 80%);
    }

    :root.light .name {
        filter: drop-shadow(-2px 2px 2px rgba(0, 0, 0, 0.5));
        color: hsl(var(--hue), 100%, 60%);
    }

    .name:focus {
        outline: none;
    }

    .buttons-bar {
        width: 100%;
        display: flex;
        justify-content: space-between;
        height: fit-content;
    }

    .buttons-bar button {
        z-index: 2;
    }

    .buttons-bar button {
        transition: var(--default-transition);
        aspect-ratio: 1 / 1;
        width: var(--size-8);
        background: transparent;
    }

    :root.dark .buttons-bar button:hover {
        filter: brightness(0.8);
    }

    :root.dark .buttons-bar button:active {
        filter: brightness(0.7);
    }

    :root.light .buttons-bar button:hover {
        filter: brightness(2);
    }

    :root.light .buttons-bar button:active {
        filter: brightness(2.5);
    }

    .progress-bar {
        position: absolute;
        width: 100%;
        height: 100%;
        animation-name: progress;
        background: conic-gradient(var(--color) var(--progress), transparent var(--progress));
    }

    :root.dark .progress-bar.destruction {
        --color: rgba(200, 0, 0, 0.25);
    }

    :root.light .progress-bar.destruction {
        --color: rgba(200, 0, 0, 0.5);
    }

    .progress-bar.reset {
        --color: color-mix(in srgb, var(--pr-font-color) 15%, transparent);
    }

    .external {
        width: 100%;
        height: 100%;
        margin: 0;
    }

    .display {
        --hue: 360deg;
        width: 100%;
        height: 100%;
        display: flex;
        flex-flow: column;
        align-items: center;
        justify-content: center;
        background-color: transparent;
        --border-color: hsl(var(--hue), 25%, 40%);
    }

    :root.dark {
        --lightness: 20%
    }

    :root.light {
        --lightness: 60%
    }

    .display.active {
        background-color: rgba(0, 0, 0, 0.1);
    }

    .display.active:hover {
        background-color: rgba(0, 0, 0, 0.2);
    }

    .display.active:active {
        background-color: rgba(0, 0, 0, 0.4);
    }

    .background {
        position: absolute;
        width: 100%;
        height: 100%;
    }

    .stopwatch-icon {
        fill: hsl(var(--hue), 20%, var(--lightness));
        width: 100%;
        height: 100%;
    }

    .stopwatch-icon path {
        filter: blur(0.5px);
    }

    .stopwatch-icon .arrow {
        transform-origin: 50% 54%;
        animation-duration: 1s;
        animation-iteration-count: infinite;
        animation-name: rotate;
        animation-timing-function: linear;
        animation-play-state: paused;
    }

    .stopwatch-icon .arrow.active {
        animation-play-state: running;
    }

    @keyframes rotate {
        0% { transform: rotate(var(--initial-rotation)); }
        100% { transform: rotate(calc(var(--initial-rotation) + 360deg)); }
    }

    .time {
        width: 80%;
        font-family: var(--font-rounded-sans);
        margin-left: auto;
        margin-right: auto;
        background-color: transparent;
        z-index: 2;
    }

    .time text {
        fill: hsl(var(--hue), 100%, calc(100% - var(--lightness)));
    }
</style>