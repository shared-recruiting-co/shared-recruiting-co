/* eslint-disable */
import type { Config, Schema } from '@markdoc/markdoc';
import type { SvelteComponent } from 'svelte';

export type ComponentsMap = Record<string, Component>;
export type Component = typeof SvelteComponent;

export type MarkdocSvelteSchema<C extends Record<string, unknown> = {}> = Schema<
	C & Config,
	Component
>;
