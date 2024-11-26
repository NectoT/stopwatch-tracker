<script lang="ts">
    import { untrack, type Snippet } from "svelte";
    import { SvelteSet } from "svelte/reactivity";
    import { fade } from "svelte/transition";

    type Props = {children: Snippet[], childrenIcons: Snippet[], displayedChild?: number};
    let { children, childrenIcons, displayedChild = 0 }: Props = $props();

    let activeChildren = $state(new SvelteSet([displayedChild]));

    $effect(() => {
        displayedChild;
        untrack(() => activeChildren.add(displayedChild));
        let id = setTimeout(() => {
            activeChildren = new SvelteSet([displayedChild]);
        }, 1000)
        return () => {
            clearTimeout(id);
        }
    })
</script>

<div class="container">
    <div class="carouselle" style="grid-template-columns: {"100% ".repeat(children.length)}; transform: translateX(-{100 * displayedChild}%);">
        {#each children as child, index}
            <div class="element">
                {#if activeChildren.has(index)}
                    {@render child()}
                {/if}
            </div>
        {/each}
    </div>

    {#snippet switchButton(cssClass: string, offset: number)}
        <button class="{cssClass} filled shadowed" transition:fade={{duration: 200}} onclick={() => displayedChild += offset}>
            {@render childrenIcons[displayedChild + offset]()}
        </button>
    {/snippet}

    {#if displayedChild > 0}
        {@render switchButton("left", -1)}
    {/if}
    {#if displayedChild < children.length - 1}
        {@render switchButton("right", 1)}
    {/if}
</div>

<style>
    .container {
        position: relative;
        display: flex;
        flex-flow: column;
        width: 100%;
        height: 100%;
        min-height: 100%;
        overflow-x: hidden;
    }

    .carouselle {
        position: relative;
        flex-grow: 1;
        width: 100%;
        display: grid;
        grid-template-rows: 100%;
        transition: transform 0.8s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    }

    .container>button {
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        top: var(--size-4);
        height: var(--size-9);
        margin: auto;
        aspect-ratio: 1 / 1;
    }

    .container>button.right {
        right: 0;
        border-radius: var(--size-5) 0 0 var(--size-5);
    }

    .container>button.left {
        left: 0;
        border-radius: 0 var(--size-5) var(--size-5) 0;
    }

    .element {
        width: 100%;
    }
</style>