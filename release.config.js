module.exports = {
  // `main` publishes stable versions; `next` publishes `-rc.N` prereleases on the
  // `next` distribution channel. Release rules follow the conventionalcommits preset:
  // feat releases a minor, fix and perf a patch, breaking changes a major.
  branches: [
    { name: 'main', channel: false },
    { name: 'next', prerelease: 'rc', channel: 'next' },
  ],
  repositoryUrl: 'git@github.com:telekom/gateway-issuer-service-go.git',
  tagFormat: '${version}',
  plugins: [
    [
      '@semantic-release/commit-analyzer',
      {
        preset: 'conventionalcommits',
        releaseRules: [{ breaking: true, release: 'major' }],
      },
    ],
    'semantic-release-export-data',
    ['@semantic-release/release-notes-generator', { preset: 'conventionalcommits' }],
    '@semantic-release/github',
  ],
};
