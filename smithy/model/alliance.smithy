$version: "2.0"

namespace ennead.wartracker

/// An Alliance resource, which has an id and descibes an alliance 
resource Alliance {
    identifiers: {
        id: Id
    }
    properties: {
        tag: Tag
        warzone: Warzone
    }
    read: GetAlliance
}

@readonly
@http(method: "GET", uri: "/alliance/{id}")
operation GetAlliance {
    input := for Alliance {
        @httpLabel
        @required
        $id
    }

    output := for Alliance {
        @required
        $id

        @required
        $tag

        @required
        $warzone
    }

    errors: [
        AllianceNotFound
    ]
}

structure AllianceData {
    name: String

    @required
    date: String

    @required
    power: Integer

    @required
    giftLevel: Integer

    @required
    memberCount: Integer

    r5Id: Id
}

@pattern("^[A-Za-z0-9]+$")
string Tag

@pattern("^[A-Za-z0-9-]+$")
string Id

integer Warzone

/// An error indicating an order could not be found
@httpError(404)
@error("client")
structure AllianceNotFound {
    message: String
    id: Id
}