import markdoc from '@markdoc/markdoc';
import type { Config, Node, RenderableTreeNode } from '@markdoc/markdoc';

const headingID = (children: RenderableTreeNode[], attributes: Record<string, unknown>) => {
	if (attributes.id && typeof attributes.id === 'string') {
		return attributes.id;
	}

	return children
		.filter((child) => typeof child === 'string')
		.join(' ')
		.replace(/[?]/g, '')
		.replace(/\s+/g, '-')
		.toLowerCase();
};

export const heading = {
	children: ['inline'],
	attributes: {
		id: { type: String },
		level: { type: Number, required: true, default: 1 }
	},
	transform(node: Node, config: Config) {
		const attributes = node.transformAttributes(config);
		const children = node.transformChildren(config);

		const id = headingID(children, attributes);

		return new markdoc.Tag(`h${node.attributes['level']}`, { ...attributes, id }, children);
	}
};

export default {
	heading
};
