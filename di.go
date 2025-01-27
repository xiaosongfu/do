package do

var DefaultInjector = New()

func getInjectorOrDefault(i *Injector) *Injector {
	if i != nil {
		return i
	}

	return DefaultInjector
}

func Provide[T any](i *Injector, provider Provider[T]) {
	name := generateServiceName[T]()

	service := newServiceLazy(name, provider)
	getInjectorOrDefault(i).set(name, service)
}

func ProvideNamed[T any](i *Injector, name string, provider Provider[T]) {
	service := newServiceLazy(name, provider)
	getInjectorOrDefault(i).set(name, service)
}

func ProvideValue[T any](i *Injector, value T) {
	name := generateServiceName[T]()

	service := newServiceEager(name, value)
	getInjectorOrDefault(i).set(name, service)
}

func ProvideNamedValue[T any](i *Injector, name string, value T) {
	service := newServiceEager(name, value)
	getInjectorOrDefault(i).set(name, service)
}

func Invoke[T any](i *Injector) (T, error) {
	name := generateServiceName[T]()
	return InvokeNamed[T](i, name)
}

func MustInvoke[T any](i *Injector) T {
	s, err := Invoke[T](i)
	must(err)
	return s
}

func InvokeNamed[T any](i *Injector, name string) (T, error) {
	serviceAny, ok := getInjectorOrDefault(i).get(name)
	if !ok {
		return empty[T](), getInjectorOrDefault(i).serviceNotFound(name)
	}

	service, ok := serviceAny.(Service[T])
	if !ok {
		return empty[T](), getInjectorOrDefault(i).serviceNotFound(name)
	}

	instance, err := service.getInstance(i)
	if err != nil {
		return empty[T](), err
	}

	getInjectorOrDefault(i).onServiceInvoke(name)

	return instance, err
}

func MustInvokeNamed[T any](i *Injector, name string) T {
	s, err := InvokeNamed[T](i, name)
	must(err)
	return s
}

func HealthCheck[T any](i *Injector) error {
	name := generateServiceName[T]()
	return getInjectorOrDefault(i).healthcheckImplem(name)
}

func HealthCheckNamed(i *Injector, name string) error {
	return getInjectorOrDefault(i).healthcheckImplem(name)
}

func Shutdown[T any](i *Injector) error {
	name := generateServiceName[T]()
	return getInjectorOrDefault(i).shutdownImplem(name)
}

func MustShutdown[T any](i *Injector) {
	name := generateServiceName[T]()
	must(getInjectorOrDefault(i).shutdownImplem(name))
}

func ShutdownNamed(i *Injector, name string) error {
	return getInjectorOrDefault(i).shutdownImplem(name)
}

func MustShutdownNamed(i *Injector, name string) {
	must(getInjectorOrDefault(i).shutdownImplem(name))
}
