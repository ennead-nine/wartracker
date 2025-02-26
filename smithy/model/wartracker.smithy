$version: "2"
namespace ennead.wartracker

use aws.protocols#restJson1
use smithy.framework#ValidationException

/// Allows users to retrieve a menu, create a coffee order, and
/// and to view the status of their orders
@title("Coffee Shop Service")
@restJson1
service Wartracker {
    version: "2024-12-23"
    resources: [
        Alliance
    ]
    errors: [
        ValidationException
    ]
}
