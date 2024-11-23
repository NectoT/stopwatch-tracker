<script lang="ts">
    import { onMount, untrack } from "svelte";
    import { AddStopwatch, DeleteStopwatch, GetStopwatch, HasStopwatch, UpdateStopwatchTime } from "@wails/go/main/App";
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
    let active = $state(false);
    let hue = $state(Math.random() * 360);

    let initialArrowRotation = $state(0);

    let loading = $state(true);

    let destructionProgress = tweened(0);
    let resetProgress = tweened(0);

    $effect(() => {
        if (!loading) {
            onLoaded();
        }
    })

    onMount(async () => {
        if (await HasStopwatch(id)) {
            let data = await GetStopwatch(id);
            hue = data.Hue;
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
            await AddStopwatch(id, "", hue)
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

    function formatTime(time: number): string {
        let date = new Date(time);
        return (
            String(date.getMinutes()).padStart(2, "0") + 
            ':' +
            String(date.getSeconds()).padStart(2, "0")
        );
    }

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
</script>

<div class="external">
    <div class="buttons-bar">
        <button onmousedown={startDestruction} onmouseup={cancelDestruction}>
            <Trash style="width: 100%; height: 100%"></Trash>
        </button>

        <button onmousedown={startReset} onmouseup={cancelReset}>
            <Reset style="width: 100%; height: 100%"></Reset>
        </button>
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
        <svg class="time" viewBox="0 -8 25 11">
            {#key time}
            <text font-size="10px" x="50%" y="0" text-anchor="middle">{formatTime(time)}</text>
            {/key}
        </svg>
    </button>
</div>

<style>
    .buttons-bar {
        position: absolute;
        width: 100%;
        display: flex;
        justify-content: space-between;
        height: fit-content;
    }

    .buttons-bar * {
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

    @property --degrees {
        syntax: '<angle>';
        initial-value: 180deg;
        inherits: false;
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

    /* :root.light .progress-bar.reset {
        --color: rgba(0, 0, 0, 0.2);
    } */

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
        background-color: transparent;
        --border-color: hsl(var(--hue), 25%, 40%);
        /* --initial-rotation: 0.542343289deg; */
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
        margin: auto;
        background-color: transparent;
        z-index: 2;
    }

    .time text {
        fill: hsl(var(--hue), 100%, calc(100% - var(--lightness)));
    }
</style>