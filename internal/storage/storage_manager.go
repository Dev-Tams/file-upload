package storage

type StorageManager struct {
    providers map[string]StorageProvider
}

func NewStorageManager() *StorageManager {
    return &StorageManager{
        providers: map[string]StorageProvider{},
    }
}

func (m *StorageManager) Register(name string, provider StorageProvider) {
    m.providers[name] = provider
}

func (m *StorageManager) Get(name string) StorageProvider {
    return m.providers[name]
}
