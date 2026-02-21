import test from 'node:test';
import assert from 'node:assert/strict';
import { appName } from './index.js';

test('appName returns frontend identifier', () => {
  assert.equal(appName(), 'GourmetGuide frontend');
});
