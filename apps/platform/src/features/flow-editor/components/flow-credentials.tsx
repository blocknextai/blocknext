import { useEffect, useState } from 'react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { bc } from '@/lib/broadcast'
import { useTranslation } from 'react-i18next'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardTitle,
} from '@/components/ui/card'

const FlowCredentials = ({
  schemas,
  secrets,
  getCredentials,
  usedCredentials,
  onSaveCredential,
  onAuthorizeOAuth2,
}) => {
  const [credentialsArray, setCredentialsArray] = useState()
  const [secretArray, setSecretArray] = useState()
  const [cSelected, setCSelected] = useState()
  const [open, setOpen] = useState(false)

  const { t } = useTranslation()

  const formatSecrets = () => {
    const haveSecrets = {}
    const mustSecrets = []
    const nSecrets = Object.assign([], secrets)
    for (let i = 0; i < nSecrets.length; i++) {
      const k = nSecrets[i]
      const s = schemas.findIndex((schema) => schema.id === k.key)
      if (s !== -1) {
        const se = schemas[s]
        if (haveSecrets[k.key] === undefined) {
          haveSecrets[k.key] = {
            id: se.id,
            name: se.name,
            description: se.description,
            isOAuth1: se.isOAuth1,
            isOAuth2: se.isOAuth2,
            items: [],
          }
        }
        haveSecrets[k.key].items.push({
          id: k.id,
          key: k.key,
          label: k.title,
          value: k.uiKey,
        })
      }
    }
    for (let i = 0; i < schemas.length; i++) {
      const s = schemas[i]
      if (haveSecrets[s.id] === undefined) {
        mustSecrets.push(s)
      }
    }
    setSecretArray(Object.values(haveSecrets))
    setCredentialsArray(mustSecrets)
  }

  // Listen for OAuth2 callback completion
  useEffect(() => {
    const handleOAuthMessage = (e) => {
      if (e.data === 'hideAuth') {
        setOpen(false)
        getCredentials().then(() => formatSecrets())
      }
    }
    bc.addEventListener('message', handleOAuthMessage)
    return () => {
      bc.removeEventListener('message', handleOAuthMessage)
    }
  }, [])

  const saveCredentials = async () => {
    const ca = cSelected
    const d = {
      title: ca.label,
      key: ca.id,
      data: ca.data,
    }
    const res = await onSaveCredential(d)
    const nsa = Object.assign({}, usedCredentials.cred)
    nsa[ca.id] = res
    usedCredentials.use(nsa)
    await getCredentials()
    formatSecrets()
    setOpen(false)
  }

  const authorizeOAuth2 = async () => {
    const ca = cSelected
    const d = {
      title: ca.label || t(ca.name),
      key: ca.id,
      data: ca.data || {},
    }

    try {
      await onAuthorizeOAuth2(d)
    } catch (error) {
      console.error('Error during OAuth2 authorization:', error)
    }
  }
  const createSecret = (item) => {
    const newItem = { ...item, label: t(item.name) }
    setCSelected(newItem)
    setOpen(true)
  }
  const updateSelectedSecret = (s, o = false) => {
    const nsa = Object.assign({}, cSelected)
    if (o) {
      if (!nsa.data) {
        nsa.data = {}
      }
      nsa.data[s.key] = s.value
    } else {
      nsa[s.key] = s.value
    }
    setCSelected(nsa)
  }
  const selectCredential = ({ settings, id }) => {
    const cred = secrets.find((secret) => secret.id === id)
    const nsa = Object.assign({}, usedCredentials.cred)
    nsa[settings.id] = cred
    usedCredentials.use(nsa)
  }
  const renderCFields = () => {
    const fields = [
      <div key="title" className="grid gap-3">
        <Label htmlFor="s-t">{t('ui.text.title')}</Label>
        <Input
          id="s-t"
          name="s-t"
          placeholder={t('ui.text.myApiKey')}
          defaultValue={cSelected?.label || ''}
          onChange={(e) =>
            updateSelectedSecret({
              key: 'label',
              value: e.target.value,
            })
          }
        />
      </div>,
    ]
    for (let i = 0; i < cSelected?.schema?.length; i++) {
      const schema = cSelected.schema[i]
      if (schema.type && schema.type !== 'hidden') {
        fields.push(
          <div key={schema.key} className="grid gap-3">
            <Label htmlFor={schema.key}>{t(schema.name)}</Label>
            <Input
              id={schema.key}
              name={schema.key}
              defaultValue={schema.default}
              disabled={schema.typeOptions?.readOnly ? 'disabled' : undefined}
              type={schema.typeOptions?.password ? 'password' : 'text'}
              placeholder={t(schema.placeholder)}
              onChange={(e) =>
                updateSelectedSecret(
                  { key: schema.key, value: e.target.value },
                  true,
                )
              }
            />
          </div>,
        )
      }
    }
    return fields
  }
  const renderMustHaves = () => {
    const items = []
    for (let i = 0; i < credentialsArray?.length; i++) {
      const ca = credentialsArray[i]
      items.push(
        <Card key={ca.id} className="w-full border-0! bg-input/30">
          <CardContent className="px-4 flex flex-col gap-4">
            <CardTitle>{t(ca.name)}</CardTitle>
            <CardDescription>{t(ca.description)}</CardDescription>
          </CardContent>
          <CardFooter className="px-4 w-full">
            <Button className="w-full" onClick={() => createSecret(ca)}>
              {t('ui.text.create')}
            </Button>
          </CardFooter>
        </Card>,
      )
    }
    return items
  }

  const renderHaves = () => {
    const items = []
    for (let i = 0; i < secretArray?.length; i++) {
      const ca = secretArray[i]
      const elems = []
      for (let j = 0; j < ca.items.length; j++) {
        const c = ca.items[j]
        elems.push(
          <SelectItem key={c.id} value={c.id}>
            {c.label}
          </SelectItem>,
        )
      }
      items.push(
        <Card key={ca.id} className="w-full border-0! bg-input/30">
          <CardContent className="px-4 flex flex-col gap-4">
            <CardTitle>{t(ca.name)}</CardTitle>
            <CardDescription>{t(ca.description)}</CardDescription>
          </CardContent>
          <CardFooter className="px-4 w-full">
            <Select
              defaultValue={
                usedCredentials.cred
                  ? usedCredentials.cred[ca.id]?.id
                  : undefined
              }
              onValueChange={(e) => selectCredential({ settings: ca, id: e })}
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder={t('ui.text.select')} />
              </SelectTrigger>
              <SelectContent>{elems}</SelectContent>
            </Select>
          </CardFooter>
        </Card>,
      )
    }
    return items
  }
  useEffect(() => {
    formatSecrets()
  }, [])
  return (
    <div className="w-full flex flex-col gap-6">
      <div className="grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 w-full grid gap-4">
        {renderHaves()}
        {renderMustHaves()}
      </div>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{t(cSelected?.name)}</DialogTitle>
            <DialogDescription>{t(cSelected?.description)}</DialogDescription>
          </DialogHeader>

          {renderCFields()}

          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">{t('ui.text.cancel')}</Button>
            </DialogClose>
            {cSelected?.isOAuth2 ? (
              <>
                <DialogClose asChild>
                  <Button variant="outline" onClick={saveCredentials}>
                    {t('ui.text.save')}
                  </Button>
                </DialogClose>
                <Button
                  variant="secondary"
                  className="bg-primary/10 hover:bg-primary/20 text-primary"
                  onClick={authorizeOAuth2}
                >
                  {t('ui.text.authorize')}
                </Button>
              </>
            ) : (
              <DialogClose asChild>
                <Button onClick={saveCredentials}>{t('ui.text.save')}</Button>
              </DialogClose>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

FlowCredentials.displayName = 'FlowCredentials'

export default FlowCredentials
