<script lang="ts" module>
    let themeState: ColorTheme = $state(ColorTheme.Light);

    export let theme = {
        get state(): ColorTheme {
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
    } from "@wails/go_timer_tracker/app";
    import { fade, fly } from "svelte/transition";
    import { ColorTheme } from "@wails/go_timer_tracker/models";
    import Plus from "./assets/icons/Plus.svelte";
    import { flip } from "svelte/animate";
    import Stopwatch from "./components/Stopwatch.svelte";
    import {SvelteSet} from "svelte/reactivity"
    import StopwatchIcon from "@assets/icons/Stopwatch.svelte"
    import StatsIcon from "@assets/icons/Stats.svelte";
    import Carouselle from "@components/Carouselle.svelte";
    import Stats from "@components/Stats.svelte";

    let timerIds: SvelteSet<string> = $state(new SvelteSet());
    
    async function toggleColorScheme() {
        let html = document.getElementsByTagName("html")[0];
        if (themeState == ColorTheme.Dark) {
            html.classList.replace("dark", "light");
            await tick();
            themeState = ColorTheme.Light;
        } else {
            html.classList.replace("light", "dark");
            await tick();
            themeState = ColorTheme.Dark;
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
        {#if themeState == ColorTheme.Light}
            <Sun></Sun>
        {:else}
            <Moon></Moon>
        {/if}
    </button>
</header>

<div class="misc">
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
</div>


{#snippet timers()}    
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
{/snippet}

{#snippet timersLogo()}
    <StopwatchIcon style="width: 75%"></StopwatchIcon>
{/snippet}

{#snippet stats()}
    <Stats></Stats>
{/snippet}

{#snippet statsLogo()}
    <StatsIcon style="width: 75%"></StatsIcon>
{/snippet}

<div style="flex-grow: 1">
    <Carouselle children={[timers, stats]} childrenIcons={[timersLogo, statsLogo]}/>
</div>

<style lang="postcss">
    header {
        display: flex;
        flex-flow: row;
        padding: var(--size-2);
    }

    .misc {
        margin: var(--size-2);
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
