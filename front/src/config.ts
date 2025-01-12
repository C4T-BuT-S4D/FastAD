const apiURL = import.meta.env.DEV
  ? 'http://127.0.0.1:8001/api'
  : window.location.origin + '/api';

const centrifugeWSURL = import.meta.env.DEV
  ? 'ws://127.0.0.1:8001/centrifuge/websocket'
  : (window.location.protocol === 'https:' ? 'wss://' : 'ws://') +
  window.location.host +
  '/centrifuge/websocket';

export { apiURL, centrifugeWSURL };
