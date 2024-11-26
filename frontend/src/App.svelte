<script lang="ts" module>
    let themeState: main.ColorTheme = $state(main.ColorTheme.Light);

    export let theme = {
        get state(): main.ColorTheme {
            return themeState;
        },
    };
</script>

<script lang="ts">
    import Moon from "./assets/icons/Moon.svelte";
    import Sun from "./assets/icons/Sun.svelte";
    import { onMount, tick } from "svelte";
    import {
        GetStopwatchIds,
        GetUserColorTheme,
        SetUserColorTheme,
    } from "@wails/go/main/App";
    import { fade, fly } from "svelte/transition";
    import { main } from "@wails/go/models";
    import Plus from "./assets/icons/Plus.svelte";
    import { flip } from "svelte/animate";
    import Stopwatch from "./components/Stopwatch.svelte";
    import {SvelteSet} from "svelte/reactivity"

    let timerIds: SvelteSet<string> = $state(new SvelteSet());
    
    async function toggleColorScheme() {
        let html = document.getElementsByTagName("html")[0];
        if (themeState == main.ColorTheme.Dark) {
            html.classList.replace("dark", "light");
            await tick();
            themeState = main.ColorTheme.Light;
        } else {
            html.classList.replace("light", "dark");
            await tick();
            themeState = main.ColorTheme.Dark;
        }
        SetUserColorTheme(themeState);
    }

    onMount(() => {
        GetUserColorTheme().then((e) => {
            themeState = e;
            document.getElementsByTagName("html")[0].classList.add(e);
        });
        GetStopwatchIds().then(ids => timerIds = new SvelteSet(ids));
    });

    async function createTimer() {
        timerIds.add(crypto.randomUUID());
        timerIds = timerIds
    }

    async function onTimerDeleted(id: string) {
        timerIds.delete(id);
    }  
</script>

<header>
    <button class="light-dark-toggle" onclick={toggleColorScheme}>
        {#if themeState == main.ColorTheme.Light}
            <Sun></Sun>
        {:else}
            <Moon></Moon>
        {/if}
    </button>
</header>

<h1>Hello, world!</h1>

{#key theme.state}
    <!-- Transitions happen because of the key changing, not because of #if directive -->
    <p in:fade={{ delay: 200, duration: 200 }} out:fade={{ duration: 200 }}>
        {#if theme.state == "light"}
            It's day time right now!
        {:else}
            It's night time! 🎃
        {/if}
    </p>
{/key}

<div class="timers-container">
    {#each timerIds as id (id)}
        <div class="item" in:fly={{duration: 200, y: 0, x: 200}} animate:flip={{duration: 200}}>
            <Stopwatch id={id} onDelete={() => onTimerDeleted(id)}></Stopwatch>
        </div>
    {/each}
    <button class="item outlined new-timer-button" onclick={e => createTimer()}>
        <Plus width="70%" height="70%"></Plus>
    </button>
</div>

<style lang="postcss">
    header {
        display: flex;
        flex-flow: row;
    }

    .item {
        position: relative;
        aspect-ratio: 1 / 1;
        max-width: var(--size-15);
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .timers-container {
        display: grid;
        justify-content: center;
        align-items: center;
        /* grid-auto-flow: column; */
        grid-template-columns: repeat(auto-fit, minmax(var(--size-13), 1fr));
        /* grid-auto-columns: minmax(var(--size-12), 100%); */
        max-width: 80%;
        align-items: center;
        justify-items: center;
        margin: auto;
        gap: var(--size-8);
    }

    .timers-container>* {
        flex-grow: 1;
        flex-basis: 0px;
    }

    .light-dark-toggle {
        flex-grow: 1;
        margin-left: auto;
        max-width: var(--size-8);
        font-size: var(--size-8);
        transition: var(--default-transition);
        background-color: transparent;
    }

    :root.light .light-dark-toggle {
        color: var(--stone-12);
    }

    :root.light .light-dark-toggle:hover {
        color: var(--stone-8);
    }

    :root.dark .light-dark-toggle {
        color: var(--violet-0);
    }

    :root.dark .light-dark-toggle:hover {
        color: var(--violet-2);
    }
</style>
