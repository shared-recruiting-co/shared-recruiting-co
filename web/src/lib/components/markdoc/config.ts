import type { Config, Node } from '@markdoc/markdoc';
import yaml from 'js-yaml';

import nodes from './nodes';

const config = (ast: Node): Config => {
	const frontmatter = ast.attributes.frontmatter ? yaml.load(ast.attributes.frontmatter) : {};

	return {
		nodes,
		tags: {
			callout: {
				render: 'Callout',
				attributes: {
					type: {
						type: String,
						default: 'note',
						matches: ['caution', 'check', 'note', 'warning']
					},
					title: {
						type: String
					}
				}
			},
			emailLabel: {
				render: 'EmailLabel',
				attributes: {
					color: {
						type: String,
						default: 'blue',
						matches: ['blue', 'green', 'red']
					},
					title: {
						type: String
					}
				}
			}
		},
		variables: {
			frontmatter
		}
	};
};

export default config;
