import { List, Paragraph, Section } from "../components/legal-content";
import LegalLayout from "../layout/legal-layout";

export default function PrivacyPage() {
  return (
    <LegalLayout title="Privacy Policy" lastUpdated="April 18, 2026">
      <Paragraph>
        ProjectMate respects your privacy. This Privacy Policy explains how we
        collect, use, and protect your information.
      </Paragraph>
      <Section title="1. Information We Collect">
        <Paragraph>We collect:</Paragraph>
        <List
          items={[
            "Name (if provided)",
            "Email address",
            "Account credentials (securely stored)",
            "Project and task data you create",
          ]}
        />
      </Section>
      <Section title="2. How We Use Your Information">
        <Paragraph>We use your data to:</Paragraph>
        <List
          items={[
            "Provide and maintain the service",
            "Authenticate users",
            "Send essential emails (verification, password reset)",
            "Improve application functionality",
          ]}
        />
      </Section>
      <Section title="3. Data Storage">
        <Paragraph>
          Your data is stored securely in our database. We take reasonable
          measures to protect your information but cannot guarantee absolute
          security.
        </Paragraph>
      </Section>
      <Section title="4. Sharing of Data">
        <Paragraph>
          We do <strong>not</strong> sell or share your personal data with third
          parties, except:
        </Paragraph>
        <List
          items={[
            "When required by law",
            "To protect the security and integrity of the service",
          ]}
        />
      </Section>
      <Section title="5. Email Communication">
        <Paragraph>We send emails for:</Paragraph>
        <List
          items={[
            "Account verification",
            "Password reset",
            "Important service updates",
          ]}
        />
      </Section>
      <Section title="6. Data Retention">
        <Paragraph>
          We retain your data as long as your account exists. You may request
          deletion of your account and associated data.
        </Paragraph>
      </Section>
      <Section title="7. Your Rights">
        <Paragraph>You may:</Paragraph>
        <List
          items={[
            "Access your data",
            "Request correction or deletion",
            "Contact us for any privacy-related concerns",
          ]}
        />
      </Section>
      <Section title="8. Changes to This Policy">
        <Paragraph>
          We may update this Privacy Policy from time to time.
        </Paragraph>
      </Section>
      <Section title="9. Contact">
        <Paragraph>
          For privacy-related questions:{" "}
          <a href="mailto:hello@contact.itsdeployedbyme.dpdns.org">
            hello@contact.itsdeployedbyme.dpdns.org
          </a>
        </Paragraph>
      </Section>
    </LegalLayout>
  );
}
