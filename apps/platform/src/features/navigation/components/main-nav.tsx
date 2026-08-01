import { ChevronRight } from 'lucide-react'
import { Link, useNavigate, useLocation } from 'react-router'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from '@/components/ui/sidebar'
import { useOrganizationStore } from '@/stores/organization'

const NavMain = ({ items, title }) => {
  const organizationId = useOrganizationStore((s) => s.organizationId)
  const navigate = useNavigate()
  const location = useLocation()

  const getUrl = (url: string) => {
    if (organizationId) {
      return url.replaceAll(':organizationId', organizationId)
    }
    return url
  }

  const isLinkActive = (url: string) => {
    if (!url || url === '#') return false
    const resolved = getUrl(url).split('?')[0]
    return location.pathname === resolved
  }
  const handleItemClick = (item, e) => {
    if (item.onClick === 'setTab' && item.tabValue) {
      e.preventDefault()
      navigate(getUrl(item.url))
    }
  }

  const renderMenu = () => {
    return (
      <SidebarGroup>
        {title && <SidebarGroupLabel>{title}</SidebarGroupLabel>}
        <SidebarMenu>
          {items.map((item, index) => {
            const hasActiveChild = item.items?.some((sub) =>
              isLinkActive(sub.url),
            )
            const itemActive = !item.items && isLinkActive(item.url)

            const isDisabled = !item.items && (!item.url || item.url === '#')

            return (
              <Collapsible
                key={index}
                asChild
                defaultOpen={item.isActive || hasActiveChild}
                open={hasActiveChild ? true : undefined}
                className="group/collapsible"
              >
                <SidebarMenuItem>
                  <CollapsibleTrigger asChild>
                    {isDisabled ? (
                      <div className={item.className}>
                        <SidebarMenuButton
                          className="cursor-default opacity-50 pointer-events-none"
                          tooltip={
                            item.tooltip ??
                            (typeof item.title === 'string'
                              ? item.title
                              : undefined)
                          }
                        >
                          {item.icon && <item.icon />}
                          {item.title}
                        </SidebarMenuButton>
                      </div>
                    ) : (
                      <Link
                        to={
                          item.url && item.url !== '#' ? getUrl(item.url) : '#'
                        }
                        className={item.className}
                        onClick={(e) => handleItemClick(item, e)}
                      >
                        <SidebarMenuButton
                          isActive={itemActive}
                          className="cursor-pointer"
                          tooltip={
                            item.tooltip ??
                            (typeof item.title === 'string'
                              ? item.title
                              : undefined)
                          }
                        >
                          {item.icon && <item.icon />}
                          {item.title}
                          {item.items && (
                            <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                          )}
                        </SidebarMenuButton>
                      </Link>
                    )}
                  </CollapsibleTrigger>
                  {item.items && (
                    <CollapsibleContent>
                      <SidebarMenuSub>
                        {item.items.map((subItem) => (
                          <SidebarMenuSubItem key={subItem.title}>
                            <SidebarMenuSubButton
                              asChild
                              isActive={isLinkActive(subItem.url)}
                            >
                              <Link
                                to={getUrl(subItem.url)}
                                className={subItem.className}
                                onClick={(e) => handleItemClick(subItem, e)}
                              >
                                {subItem.title}
                              </Link>
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                        ))}
                      </SidebarMenuSub>
                    </CollapsibleContent>
                  )}
                </SidebarMenuItem>
              </Collapsible>
            )
          })}
        </SidebarMenu>
      </SidebarGroup>
    )
  }
  return renderMenu()
}

NavMain.displayName = 'NavMain'

export { NavMain }
