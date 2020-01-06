// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.

export const environment = {
  production: false,
  API_AUTH : '/auth/users',
  API_AUTH_TOKENS: '/auth/tokens',
  API_ACCESS_CONTROL: '/access-control/access-requests',
  API_SUBSCRIPTIONS: '/subscriptions/subscriptions',
  API_STREAMS: '/streams/streams',
  API_TOKENS: '/transactions/tokens',
  API_CONTRACTS: '/transactions/contracts',
  API_EXECUTIONS: '/executions/executions',

  KUBEFLOW_URL: 'https://ar.k9s.datapace.io/_/pipeline-dashboard'
};
