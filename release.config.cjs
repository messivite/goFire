const fs = require('fs');
const path = require('path');

function updateGoVersion() {
  return {
    prepare: async (_pluginConfig, { nextRelease }) => {
      const versionFile = path.join(process.cwd(), 'config/version.go');
      let content = fs.readFileSync(versionFile, 'utf8');
      content = content.replace(/Version = "[^"]*"/, `Version = "${nextRelease.version}"`);
      fs.writeFileSync(versionFile, content);
    },
  };
}

module.exports = {
  branches: ['main'],
  plugins: [
    ['@semantic-release/commit-analyzer', { preset: 'conventionalcommits' }],
    ['@semantic-release/release-notes-generator', { preset: 'conventionalcommits' }],
    updateGoVersion,
    ['@semantic-release/changelog', { changelogFile: 'CHANGELOG.md' }],
    ['@semantic-release/git', {
      assets: ['CHANGELOG.md', 'config/version.go'],
      message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}',
    }],
    '@semantic-release/github',
  ],
};
