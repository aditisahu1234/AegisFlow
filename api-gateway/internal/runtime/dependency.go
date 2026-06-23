/*Production systems don't treat every dependency
equally. Redis and PostgreSQL may be required for
serving traffic, while OpenTelemetry or an AI
recommendation service might be optional.
 We'll build a dependency graph so the Runtime
Manager knows which components are critical and
can compute readiness automatically.
*/

package runtime

type DependencyType int

const (
	Required DependencyType = iota
	Optional
)

type Dependency struct {
	Name string
	Type DependencyType
}
