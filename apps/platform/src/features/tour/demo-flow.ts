const STARTER_NODE_ID = '0'
const VEO_NODE_ID = '1'
const VIDEO_REFERENCE = `$veo_${VEO_NODE_ID}.video`
const HANDLE_LAYOUT = 't-b'

export const DEMO_VEO_NODE_ID = VEO_NODE_ID

export const DEMO_VIDEO_REFERENCE = VIDEO_REFERENCE

export const ONBOARDING_DEMO = 'onboarding'

export const MISSING_DEMO_NODE_ID = 'youtube_upload_video'

export const onboardingDemoFlow = {
  nodes: [
    {
      id: STARTER_NODE_ID,
      nodeId: 'system_starter',
      type: 'core',
      position: { x: 0, y: 0 },
      origin: [0.5, 0.0],
      handleLayout: HANDLE_LAYOUT,
      parameters: {},
      deletable: false,
    },
    {
      id: VEO_NODE_ID,
      nodeId: 'veo',
      type: 'core',
      position: { x: 0, y: 240 },
      origin: [0.5, 0.0],
      handleLayout: HANDLE_LAYOUT,
      parameters: {
        model: 'veo-3.1-lite-generate-preview',
        prompt:
          'A cinematic 8-second product teaser: slow dolly across a desk at golden hour, warm rim light, shallow depth of field, no on-screen text.',
        aspectRatio: '16:9',
      },
    },
    {
      id: '2',
      nodeId: 'x_publish_media_post',
      type: 'core',
      position: { x: -320, y: 520 },
      origin: [0.5, 0.0],
      handleLayout: HANDLE_LAYOUT,
      parameters: {
        mediaUrls: [VIDEO_REFERENCE],
        text: 'Fresh out of the render queue 🎬 Built with BlockNext.',
      },
    },
    {
      id: '3',
      nodeId: 'instagram_publish_reels',
      type: 'core',
      position: { x: 0, y: 520 },
      origin: [0.5, 0.0],
      handleLayout: HANDLE_LAYOUT,
      parameters: {
        videoUrl: VIDEO_REFERENCE,
        caption: 'Fresh out of the render queue 🎬 Built with BlockNext.',
        shareToFeed: true,
      },
    },
  ],
  edges: [
    {
      id: 'demo-starter-veo',
      source: STARTER_NODE_ID,
      target: VEO_NODE_ID,
      sourceHandle: 'out',
      targetHandle: 'in',
      type: 'core',
      markerEnd: 'edge-circle',
    },
    {
      id: 'demo-veo-x',
      source: VEO_NODE_ID,
      target: '2',
      sourceHandle: 'out',
      targetHandle: 'in',
      type: 'core',
      markerEnd: 'edge-circle',
    },
    {
      id: 'demo-veo-instagram',
      source: VEO_NODE_ID,
      target: '3',
      sourceHandle: 'out',
      targetHandle: 'in',
      type: 'core',
      markerEnd: 'edge-circle',
    },
  ],
}

export const MISSING_DEMO_NODE_POSITION = { x: 320, y: 520 }

export const MISSING_DEMO_NODE_CANVAS_ID = String(
  onboardingDemoFlow.nodes.length,
)

export const missingDemoNodeParameters = {
  title: 'Made with BlockNext',
  categoryId: '22',
  privacy: 'unlisted',
  videoUrl: VIDEO_REFERENCE,
}
