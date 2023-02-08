<script lang="ts">
	import type { PageData } from './$types';
	export let data: PageData;
</script>

<div>
	Company
	<div>{data.company.name}</div>
	<div>{data.company.website}</div>
</div>
<div>
	Profile
	<div>{data.profile.firstName}</div>
	<div>{data.profile.lastName}</div>
</div>
