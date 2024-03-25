# spurtCMS Core Modules

The spurtCMS Core module "pkgcore", serves as the foundation for developers seeking to architect dynamic and secure web applications in Golang. Within this module there are four essential package bundles, each contributing pivotal functionalities to user authentication, member management, and content governance. With pkgcore module, developers can construct agile, user-centric applications supported by robust authentication, membership, and content management capabilities.

## 1) Auth Package

spurtCMS “Auth Package” provide functionality for managing authentication, user roles, and permissions within the spurtCMS admin application. They enable tasks such as token generation, user authentication, role management, permission assignment, and data retrieval related to authentication and authorization processes. Here are the list of functionalities carried by functions of Auth package.
<ul>
                                        <li> Migrate or set up the necessary database tables for handling authentication-related data.</li>
                                    <li>  Generates and returns an authentication token, typically used for user authentication and authorization processes.</li>
                                    <li>  Validates the authenticity and integrity of a given authentication token.</li>
                                <li>  Check if a user login attempt is valid by verifying credentials against stored user information.</li>
                                <li>  Retrieves a list of roles available in the system.</li>
                                <li>  Fetches role information based on the provided role ID.</li>
                                <li>  Creates a new role in the system.</li>
                                <li>  Updates an existing role with new information.</li>
                                <li>  Deletes a role from the system.</li>
                                <li>  Retrieves the status of a role, such as active or inactive.</li>
                                <li>  Checks if a role with a given name already exists in the system.</li>
                                <li>  Fetches all available data related to roles in the system.</li>
                                <li>  Creates a new permission for a role or user.</li>
                                <li>  Creates or updates an existing permission in the system.</li>
                                <li>  Retrieves a list of permissions associated with a specific role ID.</li>
                                <li>  Fetches details of a permission based on the provided permission ID.</li>
</ul>

## 2) Member Access Package

Functions of spurtCMS "Member Restrict" package, collectively provide functionality for managing member access control within the application. This package used in spurtCMS admin application, gives ability to CMS admin to control access to certain parts of the CMS or its content based on user roles, permissions, or other criteria. This feature is crucial for ensuring that sensitive information, administrative tools, or premium content is only accessible to authorized memebers while maintaining the integrity and security of the CMS.

<ul>
                                    <li> Migrate or set up the necessary database tables required for handling member access control.</li>
                                    <li> Retrieves information about a specific space, which could be a section or division within the application, typically related to content organization or access control.</li>
                                    <li> Retrieves details about a specific page within the admin CMS application.</li>
                                    <li> Retrieves information about a specific memeber group or access control group within the system.</li>
                                    <li>Performs checks to verify if a page requires login authentication for access.</li>
                                    <li> Retrieves a list of content items accessible to memebers based on their access rights and permissions.</li>
                                    <li> Retrieves access control information based on the provided ID, allowing administrators to manage user access rights more granularly.</li>
                                    <li> Creates access control settings for members, specifying the level of access they have to specific resources or functionalities within the system.</li>
                                    <li> Updates existing member access control settings with new configurations or permissions.</li>
                                    <li>Deletes member access control settings, revoking access to specific resources or functionalities as needed.</li>
                                    <li> Retrieves channels along with their associated entries, related to content management or publishing.</li>
                                    <li>Retrieves the count of channels available within the CMS admin, providing insights into the overall structure and organization of content.</li>

</ul>
