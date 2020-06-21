package exceptions

/*

exceptions is a package dedicated to the handling and logging of errors and messages in the system,
as well as functionality to safely run asycronous go code in a manor that is visible and protected
from exceptions bringing down a service.

The core functionality works around a Handler object which handles the execution of asyncrounous
tasks and has the ability to consistently log across the codebase. The intended use of a handler
is all processes have access to a Handler, which acts as a common method of running ansyc tasks
and structured logging of errors.

*/
