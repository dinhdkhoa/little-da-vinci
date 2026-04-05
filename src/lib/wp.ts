interface WPGraphQLParams {
  query: string;
  variables?: object;
}

/**
 * Queries the WPGraphQL endpoint
 */
export async function wpquery({ query, variables = {} }: WPGraphQLParams) {
  const wpUrl = import.meta.env.PUBLIC_WP_URL || 'https://demo.wpgraphql.com/graphql';

  const res = await fetch(wpUrl, {
    method: "post",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      query,
      variables,
    }),
  });

  if (!res.ok) {
    console.error(res);
    return {};
  }

  const { data } = await res.json();
  return data;
}

/**
 * Optionally, fetch from the REST API instead
 */
export async function getPostsRest() {
  const restUrl = import.meta.env.PUBLIC_WP_REST_URL || 'https://demo.wpgraphql.com/wp-json/wp/v2';
  const res = await fetch(`${restUrl}/posts`);
  if (!res.ok) {
    return [];
  }
  return await res.json();
}
