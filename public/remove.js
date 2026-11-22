import { maskString, maskCaptureWithGroups } from './maskString'
import { type RegexMatchWithGroups } from './exposedSecretsPatternScanner'

describe('maskString', () => {
    // Basic masking tests
    test('should mask middle portion of a string', () => {
        expect(maskString('hello world')).toBe('hel****orld')
        expect(maskString('abcdefghijklmnop')).toBe('abcde*****klmnop')
    })

    // SSH password masking tests
    test('should properly mask sshpass password', () => {
        expect(maskString('sshpass -p mypassword123 ssh user@host')).toBe('sshpass -p ************ ssh user@host')
        expect(maskString('sshpass -p secret123')).toBe('sshpass -p ************')
    })

    // LDAP password masking tests
    test('should properly mask ldap password', () => {
        expect(maskString(`ldapsearch -x -H ldap://ldap.example.com -D "cn=read" -w "mySecret!%$1234" -b "dc=example,dc=com" '(uid=user)'`)).toBe(`ldapsearch -x -H ldap://ldap.example.com -D "cn=read" -w "************ -b "dc=example,dc=com" '(uid=user)'`)
        expect(maskString(`ldapsearch -x -H ldap://ldap.example.com -D "cn=read" -w superSecret -b "dc=example,dc=com" '(uid=user)'`)).toBe(`ldapsearch -x -H ldap://ldap.example.com -D "cn=read" -w ************ -b "dc=example,dc=com" '(uid=user)'`)
        expect(maskString(`ldaps://host.domain.com\\CN=admin\\p@ssw0rd123`)).toBe(`ldaps://host.domain.com\\CN=admin\\************`)
        expect(maskString(`LDAPS://HOST\\OU=svc\\verylongpasswordthatshouldberedacted`)).toBe(`LDAPS://HOST\\OU=svc\\************`)
    })

    // Edge cases
    test('should handle empty string', () => {
        expect(maskString('')).toBe('')
    })

    test('should handle null/undefined input', () => {
        expect(maskString(null as any)).toBe(undefined)
        expect(maskString(undefined as any)).toBe(undefined)
    })

    test('should handle very short strings', () => {
        expect(maskString('ab')).toBe('*b')
        expect(maskString('abc')).toBe('a*c')
    })

    // Scanner detected secrets tests
    test('should mask various types of secrets', () => {
        // AWS Access Key
        expect(maskString('AKIA5Q8K7L9R2T1X3B0Z')).toBe('AKIA5Q8******T1X3B0Z')
        
        // API Key
        expect(maskString('sk_live_51HV9NbAcXXXXXXXXXXXXXXXXX')).toBe('sk_live_51H***********XXXXXXXXXXXX')
        
        // GitHub Token
        expect(maskString('ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZ0123')).toBe('ghp_ABCDEFG***********STUVWXYZ0123')
        
        // JWT Token
        expect(maskString('eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U')).toBe('eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.*********************************NryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U')
        
        // Private Key
        expect(maskString('-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7VJTUt9Us8cKj\nMzEfYyjiWA4R4/M2bS1GB4t7NXp98C3SC6dVMvDuictGeurT8jNbvJZHtCSuYEvu\nNMoSfm76oqFvAp8Gy0iz5sxjZmSnXyCdPEovGhLa0VzMaQ8s+CLOyS56YyCFGeJZ\n-----END PRIVATE KEY-----')).toBe('-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7VJTUt9***************************************************************************Sfm76oqFvAp8Gy0iz5sxjZmSnXyCdPEovGhLa0VzMaQ8s+CLOyS56YyCFGeJZ\n-----END PRIVATE KEY-----')
    })

    // Long string tests
    test('should handle long strings', () => {
        const longString = 'a'.repeat(100)
        const masked = maskString(longString)
        expect(masked.length).toBe(100)
        expect(masked).toMatch(/^a+[*]+a+$/)
    })
})

describe('maskCaptureWithGroups', () => {
    // Whitelist/blacklist mechanism: any group key containing "prefix" is left
    // untouched (whitelist); every other group's value gets masked (blacklist).
    test('keeps a prefix-named group visible and masks the secret-named group', () => {
        const capture = 'tokenPrefix_secretvalue1234567890'
        const pureValue = {
            groups: {
                tokenPrefix: 'tokenPrefix_',
                secretPart: 'secretvalue1234567890',
            },
        } as unknown as RegexMatchWithGroups

        expect(maskCaptureWithGroups(capture, pureValue)).toBe('tokenPrefix_secretv*******4567890')
    })

    test('masks multiple non-prefix groups, longest first', () => {
        const capture = 'user:adminSecretUsername pass:adminSecretPassword12345'
        const pureValue = {
            groups: {
                usernamePrefix: 'user:',
                username: 'adminSecretUsername',
                passwordPrefix: 'pass:',
                password: 'adminSecretPassword12345',
            },
        } as unknown as RegexMatchWithGroups

        const result = maskCaptureWithGroups(capture, pureValue)
        expect(result).toContain('user:')
        expect(result).toContain('pass:')
        expect(result).not.toContain('adminSecretUsername')
        expect(result).not.toContain('adminSecretPassword12345')
    })

    test('falls back to maskString on the full capture when there are no named groups', () => {
        const capture = 'AKIA5Q8K7L9R2T1X3B0Z'
        expect(maskCaptureWithGroups(capture, {} as RegexMatchWithGroups)).toBe(maskString(capture))
    })

    test('falls back to maskString when groups is empty', () => {
        const capture = 'AKIA5Q8K7L9R2T1X3B0Z'
        const pureValue = { groups: {} } as unknown as RegexMatchWithGroups
        expect(maskCaptureWithGroups(capture, pureValue)).toBe(maskString(capture))
    })

    // The filter is `key.toLowerCase().includes('prefix')`, a substring check,
    // not an exact-name match -- confirm mixed-case group names are still whitelisted.
    test('whitelists a prefix group regardless of casing', () => {
        const capture = 'TokenPrefix_secretvalue1234567890'
        const pureValue = {
            groups: {
                TokenPrefix: 'TokenPrefix_',
                secretPart: 'secretvalue1234567890',
            },
        } as unknown as RegexMatchWithGroups

        const result = maskCaptureWithGroups(capture, pureValue)
        expect(result).toContain('TokenPrefix_')
        expect(result).not.toContain('secretvalue1234567890')
    })

    test('masks non-prefix groups longest-first so a substring value does not leak', () => {
        const capture = 'secretPrefix_abc123 secretPrefix_abc123XYZ'
        const pureValue = {
            groups: {
                secretPrefix: 'secretPrefix_',
                shortSecret: 'abc123',
                longSecret: 'abc123XYZ',
            },
        } as unknown as RegexMatchWithGroups

        expect(maskCaptureWithGroups(capture, pureValue)).toBe('secretPrefix_ab**23 secretPrefix_abc***XYZ')
    })
})

