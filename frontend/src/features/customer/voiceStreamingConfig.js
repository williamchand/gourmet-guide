export const DEFAULT_GEMINI_VOICE_STREAMING_CONFIG = {
  model: 'gemini-2.5-flash-native-audio-preview-12-2025',
  config: {
    response_modalities: ['AUDIO'],
    system_instruction: 'You are a helpful and friendly AI assistant.'
  },
  audio: {
    format: 'pcm16',
    channels: 1,
    send_sample_rate: 16000,
    receive_sample_rate: 24000,
    chunk_size: 1024,
    input_mime_type: 'audio/pcm',
    output_mime_type: 'audio/pcm'
  }
};

function apiBaseUrl() {
  return import.meta.env.VITE_API_BASE_URL ?? '';
}

export function buildRealtimeWebSocketURL(userID, sessionID) {
  const base = apiBaseUrl().replace(/\/$/, '');
  const wsProtocol = base.startsWith('https://') ? 'wss://' : 'ws://';
  if (base.startsWith('http://') || base.startsWith('https://')) {
    const host = base.replace(/^https?:\/\//, '');
    return `${wsProtocol}${host}/ws/${encodeURIComponent(userID)}/${encodeURIComponent(sessionID)}`;
  }
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
  return `${protocol}://${window.location.host}/ws/${encodeURIComponent(userID)}/${encodeURIComponent(sessionID)}`;
}

export async function fetchVoiceStreamingConfig() {
  const response = await fetch(`${apiBaseUrl()}/v1/realtime/voice-config`);
  if (!response.ok) {
    throw new Error('failed to load backend voice streaming config');
  }
  return response.json();
}

export function buildRealtimeAudioMessage(data, streamingConfig) {
  const inputMimeType =
    streamingConfig?.audio?.input_mime_type ??
    DEFAULT_GEMINI_VOICE_STREAMING_CONFIG.audio.input_mime_type;
  return {
    data,
    mime_type: inputMimeType
  };
}
