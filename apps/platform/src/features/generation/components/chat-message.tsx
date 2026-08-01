import Markdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Sparkles } from 'lucide-react'
import WorkflowPreview from '@/features/generation/components/workflow-preview'

function parseWorkflowBlocks(text) {
  const parts = []
  // Support both ```json:workflow and ```json formats with flexible whitespace
  const regex = /```(?:json:workflow|json)\s*([\s\S]*?)```/g
  let lastIndex = 0
  let match

  while ((match = regex.exec(text)) !== null) {
    const content = match[1].trim()
    // Only treat as workflow if it looks like a workflow object (has nodes or edges)
    const isWorkflow =
      content.includes('"nodes"') || content.includes('"edges"')

    if (match.index > lastIndex) {
      parts.push({
        type: 'text',
        content: text.slice(lastIndex, match.index),
      })
    }

    if (isWorkflow) {
      parts.push({ type: 'workflow', content })
    } else {
      // Not a workflow JSON, keep as code block
      parts.push({ type: 'text', content: match[0] })
    }
    lastIndex = match.index + match[0].length
  }

  if (lastIndex < text.length) {
    parts.push({ type: 'text', content: text.slice(lastIndex) })
  }

  return parts
}

const MarkdownContent = ({ content }) => {
  return (
    <Markdown
      remarkPlugins={[remarkGfm]}
      components={{
        h1: ({ children }) => (
          <h1 className="text-lg font-semibold mt-4 mb-1.5">{children}</h1>
        ),
        h2: ({ children }) => (
          <h2 className="text-base font-semibold mt-4 mb-1.5">{children}</h2>
        ),
        h3: ({ children }) => (
          <h3 className="text-sm font-semibold mt-4 mb-1.5">{children}</h3>
        ),
        p: ({ children }) => <p className="mb-2 last:mb-0">{children}</p>,
        ul: ({ children }) => (
          <ul className="list-disc list-inside mb-2 space-y-0.5">{children}</ul>
        ),
        ol: ({ children }) => (
          <ol className="list-decimal list-inside mb-2 space-y-0.5">
            {children}
          </ol>
        ),
        li: ({ children }) => <li>{children}</li>,
        strong: ({ children }) => (
          <strong className="font-semibold">{children}</strong>
        ),
        code: ({ inline, children }) =>
          inline ? (
            <code className="bg-muted/70 px-1.5 py-0.5 rounded text-[13px]">
              {children}
            </code>
          ) : (
            <code>{children}</code>
          ),
        pre: ({ children }) => (
          <pre className="bg-muted/70 rounded-lg p-3 my-2.5 overflow-x-auto text-[13px]">
            {children}
          </pre>
        ),
        a: ({ href, children }) => (
          <a
            href={href}
            target="_blank"
            rel="noopener noreferrer"
            className="text-primary underline underline-offset-2"
          >
            {children}
          </a>
        ),
        hr: () => <hr className="my-3 border-border" />,
        blockquote: ({ children }) => (
          <blockquote className="border-l-2 border-primary/30 pl-3 my-2 text-muted-foreground italic">
            {children}
          </blockquote>
        ),
        table: ({ children }) => (
          <div className="my-2.5 overflow-x-auto rounded-lg border border-border">
            <table className="w-full text-xs">{children}</table>
          </div>
        ),
        thead: ({ children }) => (
          <thead className="bg-muted/50">{children}</thead>
        ),
        th: ({ children }) => (
          <th className="px-3 py-1.5 text-left font-semibold border-b border-border">
            {children}
          </th>
        ),
        td: ({ children }) => (
          <td className="px-3 py-1.5 border-b border-border">{children}</td>
        ),
        tr: ({ children }) => <tr className="hover:bg-muted/30">{children}</tr>,
      }}
    >
      {content}
    </Markdown>
  )
}

const ChatMessage = ({ message, onApplyWorkflow }) => {
  const isUser = message.role === 'user'
  const parts = parseWorkflowBlocks(message.content || '')

  if (isUser) {
    return (
      <div className="flex justify-end">
        <div className="max-w-[75%] bg-primary text-primary-foreground rounded-2xl rounded-br-md px-4 py-2.5 text-sm">
          {message.content}
        </div>
      </div>
    )
  }

  return (
    <div className="flex gap-3 w-full">
      <div className="shrink-0 w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center mt-0.5">
        <Sparkles className="w-4 h-4 text-primary" />
      </div>
      <div className="flex-1 min-w-0 text-sm">
        {parts.map((part, index) => {
          if (part.type === 'workflow') {
            return (
              <WorkflowPreview
                key={index}
                json={part.content}
                onApplyWorkflow={onApplyWorkflow}
              />
            )
          }
          return <MarkdownContent key={index} content={part.content} />
        })}
      </div>
    </div>
  )
}

ChatMessage.displayName = 'ChatMessage'

export default ChatMessage
